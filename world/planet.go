package world

import (
	"fmt"
	"math"
	"math/rand"
)

type Planet struct {
	Name           string
	Location       Coordinates
	Resources      map[Resource]int
	Population     int
	DangerLevel    int
	ResourceDemand map[Resource]int
	IsHabitable    bool
	IsHarvestable  bool
}

func (w *World) IsHabitable(probability float64) bool {

	generated := w.RandomNumber.Float64()
	return generated < probability
}

func (w *World) GeneratePlanetsInZone(numPlanets int, zone Zone, zoneType ZoneType) map[string]*Planet {
	zonePlanets := make(map[string]*Planet, numPlanets)
	for i := 0; i < numPlanets; i++ {
		distanceFromZoneCenter := w.RandomNumber.Float64() * 100
		angle := w.RandomNumber.Float64() * 2 * 3.14159

		planetLocation := Coordinates{
			X: zone.CentralPoint.X + distanceFromZoneCenter*math.Cos(angle),
			Y: zone.CentralPoint.Y + distanceFromZoneCenter*math.Sin(angle),
		}

		isHabitable := w.IsHabitable(zoneType.HabitableProbability)

		population := 0

		if isHabitable {
			population = w.randomInt(zoneType.LowerPopulation, zoneType.HigherPopulation)
		}

		planet := &Planet{
			Name:          fmt.Sprintf("%s-Planet-%d", zone.Name, i+1),
			Location:      planetLocation,
			Population:    population,
			DangerLevel:   w.randomInt(zone.DangerRange[0], zone.DangerRange[1]),
			IsHabitable:   isHabitable,
			IsHarvestable: !isHabitable,
		}

		planet.Resources = GeneratePlanetResources(*w, zone, *planet)

		w.Planets[planet.Name] = planet
		zonePlanets[planet.Name] = planet
	}

	return zonePlanets
}

func GeneratePlanetResources(world World, zone Zone, planet Planet) map[Resource]int {
	resources := make(map[Resource]int, 4)

	for res := range world.AllResources {
		if shouldIncludeResource(world, res, planet) {
			resources[res] = rand.Intn(1_000_000)
			continue
		}

		resources[res] = 0
	}

	resources[zone.ResourceProfile.Primary] = resources[zone.ResourceProfile.Primary]*world.RandomNumber.Intn(5) + 1
	resources[zone.ResourceProfile.Secondary] = resources[zone.ResourceProfile.Secondary]*world.RandomNumber.Intn(3) + 1

	return resources
}
