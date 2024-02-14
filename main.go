package main

import (
	"fmt"
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
	Steps []string
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

		recipes = append(recipes, newRecipe)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.Visit("https://www.cookwell.com/recipe/jalapeno-ranch-salad-w-adobo-chicken")
}