package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Recipe struct {
	Title string
	Descripition string
	Ingredients []string
	Tags []string
	Author string
	Link string
	ImageUrl string
}

func main () {
	c := colly.NewCollector(
		colly.AllowedDomains("www.cookwell.com"),
		colly.CacheDir("./cook_box_cache"),
	)

	recipes := make([]Recipe, 0)


	c.OnHTML("section", func(e *colly.HTMLElement) {
		newRecipe := Recipe{}
		newRecipe.Title = e.ChildText("h1.text-heading-1")
		newRecipe.Descripition = e.ChildText("div.container.col-span-2.flex.flex-col.gap-8.py-8.lg\\:pr-10 > div.prose > h2")
		imageUrl := e.ChildText("img")
		log.Print("Image URL: ", imageUrl)
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
		
		e.ForEach("li[data-testid=ingredient-group]", func(_ int, el *colly.HTMLElement) {
			el.ForEach("ul > li > div.flex-grow > div > span", func(_ int, el *colly.HTMLElement) {
				subIngredient := el.Text
				newRecipe.Ingredients = append(newRecipe.Ingredients, subIngredient)
			})
		})
		recipes = append(recipes, newRecipe)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.cookwell.com/recipe/jalapeno-ranch-salad-w-adobo-chicken")

}