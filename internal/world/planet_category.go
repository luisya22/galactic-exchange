package world

import (
	"encoding/json"
	"fmt"
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

	file, err := gamedata.Files.Open("categorydata/categories.json")
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

	mainCategoryIndex := w.randomInt(0, len(categories)-1)
	secondaryCategoryIndex := w.randomInt(0, len(categories)-1)

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
// TODO: Loop per planet consuming random resources.
// TODO: It would be basics based on population.
// TODO: Also by technology
// TODO: Add bonus consumptions, this would have resource and endTime

func (w *World) simulateConsumption() {
	for range w.newDayChan {
		w.consumeResources()
	}
}

func (w *World) consumeResources() {
	for _, planet := range w.Planets {

		if !planet.IsHabitable {
			continue
		}

		cp := planet.CategoryProfile

		for resourceName, rc := range cp.mainProfile.resourceConsumption {
			w.processResourceConsumption(planet.Name, resourceName, rc.minConsumption, rc.maxConsumption)
		}

		for resourceName, rc := range cp.secondaryProfile.resourceConsumption {
			w.processResourceConsumption(planet.Name, resourceName, rc.minConsumption, rc.maxConsumption)
		}

		_, _ = w.DepletePlanetResource(planet.Name, "food", planet.Population)
		// TODO: food restock should happen only if food production * 30 < actual stock
	}
}

func (w *World) processResourceConsumption(planetId string, resourceName string, minConsumption int, maxConsumption int) {
	quantity := w.randomInt(minConsumption, maxConsumption)
	remaning, _ := w.DepletePlanetResource(planetId, resourceName, quantity)

	weeklyConsumption := quantity * 7

	w.restockResources(planetId, resourceName, weeklyConsumption, remaning)
}

func classifyResourceLevel(weeklyConsumption int, totalStorage int) string {
	if totalStorage == 0 {
		return "depleted"
	}

	consumptionPercentage := (weeklyConsumption / totalStorage) * 100

	switch {
	case consumptionPercentage < 10:
		return "high"
	case consumptionPercentage >= 10 && consumptionPercentage < 50:
		return "medium"
	case consumptionPercentage >= 50 && consumptionPercentage < 75:
		return "low"
	default:
		return "depleted"
	}
}

func purchaseProbability(resourceLevel string) float64 {
	switch resourceLevel {
	case "high":
		return 0.10
	case "medium":
		return 0.50
	case "low":
		return 0.75
	case "depleted":
		return 1.00
	default:
		return 0.00
	}
}

func (w *World) restockResources(planetId string, resource string, weeklyConsumption int, totalStorage int) {
	resourceLevel := classifyResourceLevel(weeklyConsumption, totalStorage)
	purchaseProb := purchaseProbability(resourceLevel)

	randomFloat := w.RandomNumber.Float64()

	wantsToBuy := w.randomInt(weeklyConsumption, weeklyConsumption*4)

	if randomFloat <= purchaseProb {
		// TODO: Make purchase
		fmt.Printf("%v wants to buy %v %v\n", planetId, wantsToBuy, resource)
	}
}

// TODO: planets should analyze their resource scarcity
