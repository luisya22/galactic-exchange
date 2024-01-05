package world_test

import (
	"math/rand"
	"testing"

	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/world"
)

const planet1Name = "Zone-1-Planet-1"
const resourceQuantity = 1_000_000

func createTestWorld(t *testing.T, gameChannels *gamecomm.GameChannels) *world.World {
	t.Helper()

	randomNumber := rand.New(rand.NewSource(0))
	resourceRarity := map[world.Resource]world.Rarity{
		world.Gold:  world.Common,
		world.Iron:  world.Common,
		world.Water: world.Scarce,
		world.Food:  world.Rare,
	}

	allResources := world.CreateWorldResources()
	allZoneTypes := world.CreateZoneTypes()

	w := &world.World{
		ResourceRarity: resourceRarity,
		AllResources:   allResources,
		AllZoneTypes:   allZoneTypes,
		RandomNumber:   randomNumber,
		Workers:        100,
		WorldChan:      gameChannels.WorldChannel,
		Size:           10_000,
		Zones:          make(map[string]*world.Zone),
		Planets:        make(map[string]*world.Planet),
	}

	w.LayerBoundaries = world.GenerateLayerBoundaries(w)

	coordinates := world.Coordinates{10, 10}

	zone1 := createTestZone("Zone-1", coordinates)
	w.Zones[zone1.Name] = zone1

	planetCoordinates := world.Coordinates{20, 20}
	createTestPlanet(w, zone1, planet1Name, true, planetCoordinates, 0, 1)

	return w
}

func createTestZone(zoneName string, zoneLocation world.Coordinates) *world.Zone {

	dangerLevel := 1
	resourcesProfile := world.ResourceProfile{
		Primary:   world.Gold,
		Secondary: world.Water,
	}

	zone := &world.Zone{
		Name:            zoneName,
		CentralPoint:    world.Coordinates{10, 10},
		DangerRange:     [2]int{dangerLevel, dangerLevel + 10},
		ResourceProfile: resourcesProfile,
		ZoneType:        world.SectorOne,
		Planets:         make(map[string]*world.Planet),
	}

	return zone

}

func createTestPlanet(w *world.World, z *world.Zone, planetName string, isHabitable bool, planetLocation world.Coordinates, population int, dangerLevel int) {

	planet := world.Planet{
		Name: planetName, Location: planetLocation,
		Population:    population,
		DangerLevel:   dangerLevel,
		IsHabitable:   isHabitable,
		IsHarvestable: !isHabitable,
	}

	resources := make(map[world.Resource]int)

	resources[world.Gold] = resourceQuantity
	resources[world.Iron] = resourceQuantity
	resources[world.Water] = resourceQuantity
	resources[world.Food] = resourceQuantity

	planet.Resources = resources

	z.Planets[planet.Name] = &planet
	w.Planets[planet.Name] = &planet

}
