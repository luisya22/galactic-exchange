package world

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/resource"
)

type World struct {
	Planets         map[string]*Planet
	Zones           map[string]*Zone
	AllResources    map[string]resource.Resource
	AllZoneTypes    map[LayerName]ZoneType
	LayerBoundaries []float64
	RandomNumber    *rand.Rand
	RW              sync.RWMutex
	Workers         int
	WorldChan       chan gamecomm.WorldCommand
	Size            float64
	Categories      map[string]Category
}

func New(gameChannels *gamecomm.GameChannels, resources map[string]resource.Resource) *World {

	randomnumber := rand.New(rand.NewSource(time.Now().UnixNano()))

	allZoneTypes := CreateZoneTypes()
	world := &World{
		AllResources: resources,
		AllZoneTypes: allZoneTypes,
		RandomNumber: randomnumber,
		Workers:      100,
		WorldChan:    gameChannels.WorldChannel,
		Size:         10_000,
	}

	world.Zones = make(map[string]*Zone, 1000)
	world.Planets = make(map[string]*Planet, 8000)

	world.LayerBoundaries = GenerateLayerBoundaries(world)

	world.GenerateZones(1000)

	go world.Listen()

	return world
}

func (w *World) randomInt(min, max int) int {
	return min + w.RandomNumber.Intn(max-min+1)
}

func CreateZoneTypes() map[LayerName]ZoneType {
	return map[LayerName]ZoneType{
		SectorOne: {
			Name:                 SectorOne,
			LowerDanger:          0,
			HigherDanger:         10,
			LowerPopulation:      0,
			HigherPopulation:     7_000_000_000,
			LowerPlanetsAmount:   1,
			HigherPlanetsAmount:  10,
			HabitableProbability: .75,
			Index:                0,
			MapPercentage:        .15,
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
			Index:                1,
			MapPercentage:        .20,
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
			Index:                2,
			MapPercentage:        .20,
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
			Index:                3,
			MapPercentage:        .45,
		},
	}

}

func GenerateLayerBoundaries(w *World) []float64 {
	layerBoundaries := []float64{}

	// Calculate layer boundaries
	currentBoundary := 0.0
	for _, z := range w.AllZoneTypes {
		currentBoundary += w.Size / 2 * z.MapPercentage

		layerBoundaries = append(layerBoundaries, currentBoundary)
	}

	return layerBoundaries
}

func (w *World) GetZoneByIndex(index int) (*ZoneType, error) {
	for _, z := range w.AllZoneTypes {
		if z.Index == index {
			return &z, nil
		}
	}

	return nil, fmt.Errorf("error: zone not found with inde %v", index)
}

func (w *World) GetZoneIds() []string {
	ids := []string{}
	for _, zone := range w.Zones {
		ids = append(ids, zone.Name)
	}

	return ids
}
