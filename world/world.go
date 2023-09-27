package world

import (
	"math/rand"
	"time"
)

type World struct {
	Planets        map[string]*Planet
	Zones          map[string]*Zone
	ResourceRarity map[Resource]Rarity
	AllResources   map[Resource]struct{}
	AllZoneTypes   map[LayerName]ZoneType
	RandomNumber   *rand.Rand
}

func New() *World {

	randomnumber := rand.New(rand.NewSource(time.Now().UnixNano()))

	resourceRarity := map[Resource]Rarity{
		Gold:  Common,
		Iron:  Common,
		Water: Scarce,
		Food:  Rare,
	}

	allResources := map[Resource]struct{}{
		Gold:  {},
		Iron:  {},
		Water: {},
		Food:  {},
	}

	allZoneTypes := map[LayerName]ZoneType{
		SectorOne: {
			Name:                 SectorOne,
			LowerDanger:          0,
			HigherDanger:         10,
			LowerPopulation:      0,
			HigherPopulation:     7_000_000_000,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  10,
			HabitableProbability: .75,
		},
		SectorTwo: {
			Name:                 SectorTwo,
			LowerDanger:          10,
			HigherDanger:         25,
			LowerPopulation:      0,
			HigherPopulation:     5_000_000_000,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  10,
			HabitableProbability: .50,
		},
		SectorThree: {
			Name:                 SectorThree,
			LowerDanger:          25,
			HigherDanger:         50,
			LowerPopulation:      0,
			HigherPopulation:     1_000_000,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  15,
			HabitableProbability: .20,
		},
		SectorFour: {
			Name:                 SectorFour,
			LowerDanger:          50,
			HigherDanger:         100,
			LowerPopulation:      0,
			HigherPopulation:     0,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  20,
			HabitableProbability: 0,
		},
	}

	world := &World{
		ResourceRarity: resourceRarity,
		AllResources:   allResources,
		AllZoneTypes:   allZoneTypes,
		RandomNumber:   randomnumber,
	}

	world.Zones = make(map[string]*Zone, 1000)
	world.Planets = make(map[string]*Planet, 100)

	world.GenerateZones(10_000, 1000)

	return world

}

func (w *World) randomInt(min, max int) int {
	return min + w.RandomNumber.Intn(max-min+1)
}
