package world

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/luisya22/galactic-exchange/gamecomm"
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
	RW             sync.RWMutex
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

		planet.RW.Lock()
		GeneratePlanetResources(w, zone, planet)
		planet.RW.Unlock()

		w.Planets[planet.Name] = planet
		zonePlanets[planet.Name] = planet
	}

	return zonePlanets
}

func GeneratePlanetResources(world *World, zone Zone, planet *Planet) {
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

	planet.Resources = resources
}

func (p *Planet) copy() gamecomm.Planet {
	return gamecomm.Planet{
		Name:          p.Name,
		Location:      gamecomm.Coordinates{X: p.Location.X, Y: p.Location.Y},
		Population:    p.Population,
		DangerLevel:   p.DangerLevel,
		IsHabitable:   p.IsHabitable,
		IsHarvestable: p.IsHarvestable,
	}
}

func (w *World) GetPlanet(planetId string) (gamecomm.Planet, error) {
	var planet *Planet
	var ok bool

	w.RW.RLock()
	if planet, ok = w.Planets[planetId]; !ok {
		return gamecomm.Planet{}, fmt.Errorf("Planet not found: %v", planetId)
	}

	defer w.RW.RUnlock()

	return planet.copy(), nil
}

func (w *World) getPlanetReference(planetId string) (*Planet, error) {
	var planet *Planet
	var ok bool

	if planet, ok = w.Planets[planetId]; !ok {
		return nil, fmt.Errorf("error: planet not found: %v", planetId)
	}

	return planet, nil
}

func (w *World) RemoveResourcesFromPlanet(planetId string, resourceName Resource, amount int) (int, error) {
	w.RW.Lock()
	planet, err := w.getPlanetReference(planetId)
	if err != nil {
		return 0, err
	}
	w.RW.Unlock()

	planet.RW.Lock()
	defer planet.RW.Unlock()

	resources, ok := planet.Resources[resourceName]

	if !ok || amount > resources {
		return 0, fmt.Errorf("error: not enough resources")
	}

	planet.Resources[resourceName] -= amount

	return planet.Resources[resourceName], nil
}

func (w *World) AddResourcesToPlanet(planetId string, resourceName Resource, amount int) (int, error) {

	if amount < 0 {
		return 0, fmt.Errorf("error: amount should be greater than zero")
	}

	w.RW.RLock()
	planet, err := w.getPlanetReference(planetId)
	if err != nil {
		return 0, err
	}
	w.RW.RUnlock()

	planet.RW.Lock()
	defer planet.RW.Unlock()

	if _, ok := planet.Resources[resourceName]; !ok {
		planet.Resources[resourceName] = 0
	}

	planet.Resources[resourceName] += amount

	return planet.Resources[resourceName], nil
}
