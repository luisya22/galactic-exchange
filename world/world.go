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
	AllZoneTypes   map[ZoneName]ZoneType
	RandomNumber   *rand.Rand
}

func NewWorld() *World {

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

	allZoneTypes := map[ZoneName]ZoneType{
		GuardianSectors: {
			Name:                 GuardianSectors,
			LowerDanger:          0,
			HigherDanger:         10,
			LowerPopulation:      0,
			HigherPopulation:     7_000_000_000,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  10,
			HabitableProbability: .90,
		},
		TradeLanes: {
			Name:                 TradeLanes,
			LowerDanger:          10,
			HigherDanger:         25,
			LowerPopulation:      0,
			HigherPopulation:     5_000_000_000,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  10,
			HabitableProbability: .50,
		},
		OutlawQuadrants: {
			Name:                 OutlawQuadrants,
			LowerDanger:          25,
			HigherDanger:         50,
			LowerPopulation:      0,
			HigherPopulation:     1_000_000,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  15,
			HabitableProbability: .20,
		},
		DarkMatterZones: {
			Name:                 DarkMatterZones,
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

	world.Zones = make(map[string]*Zone, 100000)
	world.Planets = make(map[string]*Planet, 100)

	world.GenerateZones(10_000, 100_000)

	return world

}

func (w *World) randomInt(min, max int) int {
	return min + w.RandomNumber.Intn(max-min+1)
}
