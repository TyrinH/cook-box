package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"database/sql"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Recipe struct {
	ID int64
	Title string
	Descripition string
	Ingredients []string
	Tags []string
	Author string
	recipeLink string
	ImageUrl string
	Steps []string
}

var db *sql.DB

func main () {
	godotenv.Load()
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	log.Print()
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", DB_USER, DB_NAME, DB_PASSWORD)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	recipe := scrapeWebsite("https://www.cookwell.com/recipe/jalapeno-ranch-salad-w-adobo-chicken")
	
	_, addRecipeErr := addRecipe(recipe)
	if addRecipeErr != nil {
		log.Fatal(err)
	}

	// foundRecipe, err := recipeById(newRecipeId)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func scrapeWebsite (url string) (Recipe) {
	splitUrls := strings.SplitAfter(url, ".com")
	baseUrl := strings.Replace(splitUrls[0], "https://",  "", 1)
	log.Print("BASEURL: ", baseUrl)


	c := colly.NewCollector(
		colly.AllowedDomains(baseUrl),
		colly.CacheDir("./cook_box_cache"),
	)

	newRecipe := Recipe{}
	newRecipe.recipeLink = url

	c.OnHTML("section", func(e *colly.HTMLElement) {

		if e.ChildText("h1.text-heading-1") != "" {
			newRecipe.Title = e.ChildText("h1.text-heading-1")
		}
		if e.ChildText("div.container.col-span-2.flex.flex-col.gap-8.py-8.lg\\:pr-10 > div.prose > h2") != "" {
			newRecipe.Descripition = e.ChildText("div.container.col-span-2.flex.flex-col.gap-8.py-8.lg\\:pr-10 > div.prose > h2")
		}
		rawTags := e.ChildText("div.flex.flex-row.items-center.gap-1")
		listOfSpans := e.ChildText("div.container.col-span-2.flex.flex-col.gap-8.py-8.lg\\:pr-10 > span")

		if strings.ContainsAny(listOfSpans, "By") {
			newRecipe.Author = strings.Trim(listOfSpans, "By")
		}
		if strings.ContainsAny(rawTags, "Tags:") {
			tags := strings.Trim(rawTags, "Tags:")
			tagsSanitized := strings.Split(tags,`,`)
			for i := 0; i < len(tagsSanitized); i++ {
				newRecipe.Tags = append(newRecipe.Tags, tagsSanitized[i])
			}
		}
			e.ForEach("ul > li > div.flex-grow > div > span", func(_ int, el *colly.HTMLElement) {
				subIngredient := el.Text
				newRecipe.Ingredients = append(newRecipe.Ingredients, subIngredient)
			})
			
			e.ForEach("div.grid.grid-cols-1.gap-8.lg\\:hidden > section.prose.flex-col > div", func(_ int, el *colly.HTMLElement) {
				step := el.Text
				newRecipe.Steps = append(newRecipe.Steps, step)
			})

	})
	
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.Visit(url)
	log.Print("New recipe Title: ", newRecipe.Title)

	return newRecipe
}

func addRecipe(rec Recipe) (int64, error) {
	var id int64
	err := db.QueryRow(`INSERT INTO recipes(title, descripition, author, websitelink) VALUES($1, $2, $3, $4) RETURNING id`, rec.Title, rec.Descripition, rec.Author, rec.recipeLink).Scan(&id)
	if err != nil {
        return 0, fmt.Errorf("addRecipe: %v", err)
    }
    if err != nil {
        return 0, fmt.Errorf("addRecipe: %v", err)
    }
    return id, nil
}

func recipeById(id int64) (Recipe, error) {
	var rec Recipe

	row := db.QueryRow("SELECT * FROM recipes WHERE id = $1", id)
	if err := row.Scan(&rec.ID, &rec.Title, &rec.Descripition, &rec.Author, &rec.recipeLink); err != nil {
        if err == sql.ErrNoRows {
            return rec, fmt.Errorf("recById %d: no such recipe", id)
        }
        return rec, fmt.Errorf("recipeById %d: %v", id, err)
    }
    return rec, nil
}