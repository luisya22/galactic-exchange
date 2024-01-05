package world_test

import (
	"testing"

	"github.com/luisya22/galactic-exchange/world"
)

func TestGeneratePlanetsInZone(t *testing.T) {
	// numPlanets int, zone Zone, zoneType ZoneType
	testMap := []struct {
		name       string
		numPlanets int
		zone       *world.Zone
		wants      []world.Planet
	}{}

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {
			// Create world
			allZoneTypes := world.CreateZoneTypes()
			w := world.World{
				Zones: map[string]*world.Zone{
					tt.zone.Name: tt.zone,
				},
				AllZoneTypes: allZoneTypes,
				Size:         1000,
			}

			zoneType := w.AllZoneTypes[tt.zone.ZoneType]

			// Generate planets
			w.GeneratePlanetsInZone(tt.numPlanets, *tt.zone, zoneType)
			// Validate that planet follow the ranges
			// Distance
			// minDistance := w.Size / 2
			// habitable probability
			// population
			// danger level

		})
	}
}

func TestGeneratePlanetResources(t *testing.T) {
	// testMap := []struct {
	// 	name   string
	// 	planet world.Planet
	// 	zone   world.Zone
	// }
}
