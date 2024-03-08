package world

import (
	"encoding/json"
	"log"

	"github.com/luisya22/galactic-exchange/internal/gamedata"
)

// planetCategories is the main categories of the planet
type planetCategories struct {
	mainProfile      categoryProfile
	secondaryProfile categoryProfile
}

type categoryProfile struct {
	category            string
	level               uint
	resourceConsumption resourceConsumption
}

type resourceConsumption map[string]consumptionValues

type consumptionValues struct {
	minConsumption int
	maxConsumption int
}

type Category struct {
	Name        string
	Description string
	Resources   resourceSlice
}

type resourceSlice []string

func (rs *resourceSlice) Contains(s string) bool {
	for _, r := range *rs {
		if s == r {
			return true
		}
	}

	return false
}

func loadCategories() map[string]Category {
	categories := make(map[string]Category, 6)

	file, err := gamedata.Files.Open("categoriesdata/categories.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&categories)
	if err != nil {
		log.Fatal(err.Error())
	}

	return categories
}

func (w *World) generatePlanetCategoryProfile() planetCategories {

	categories := []string{}
	for categoryId := range w.Categories {
		categories = append(categories, categoryId)
	}

	mainCategoryIndex := w.randomInt(0, len(categories))
	secondaryCategoryIndex := w.randomInt(0, len(categories))

	mainCategory := categories[mainCategoryIndex]
	secondaryCategory := categories[secondaryCategoryIndex]

	// Randomly select level
	mainLevel := w.randomInt(0, 99)
	secondaryLevel := w.randomInt(0, 20)

	// Depending on level select consumptionValues

	mainConsumption := make(resourceConsumption)
	secondaryConsumption := make(resourceConsumption)

	for resourceName := range w.AllResources {

		mc := w.Categories[mainCategory]

		if mc.Resources.Contains(resourceName) {
			mainMinConsumption := w.randomInt(100*mainLevel, 200*mainLevel)
			mainMaxConsumption := w.randomInt(200*mainLevel, 300*mainLevel)

			mainConsumption[resourceName] = consumptionValues{
				minConsumption: mainMinConsumption,
				maxConsumption: mainMaxConsumption,
			}
		} else {
			sc := w.Categories[secondaryCategory]

			if sc.Resources.Contains(resourceName) {
				secondaryMinConsumption := w.randomInt(10*secondaryLevel, 20*secondaryLevel)
				secondaryMaxConsumption := w.randomInt(20*secondaryLevel, 30*secondaryLevel)

				secondaryConsumption[resourceName] = consumptionValues{
					minConsumption: secondaryMinConsumption,
					maxConsumption: secondaryMaxConsumption,
				}
			}

		}

	}

	mainCategoryProfile := categoryProfile{
		category:            mainCategory,
		level:               uint(mainLevel),
		resourceConsumption: mainConsumption,
	}

	secondaryCategoryProfile := categoryProfile{
		category:            secondaryCategory,
		level:               uint(secondaryLevel),
		resourceConsumption: secondaryConsumption,
	}

	return planetCategories{
		mainProfile:      mainCategoryProfile,
		secondaryProfile: secondaryCategoryProfile,
	}
}

// Each level would give x amount of more consumption (also based on population)
