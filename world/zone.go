package world

import (
	"fmt"
	"math"
)

type ZoneType struct {
	Name                 LayerName
	LowerDanger          int
	HigherDanger         int
	LowerPopulation      int
	HigherPopulation     int
	LowerPlanetsAmount   int
	HigherPlanetsAmount  int
	HabitableProbability float64
	Index                int
	MapPercentage        float64
}

type Zone struct {
	Name            string
	CentralPoint    Coordinates
	DangerRange     [2]int
	ResourceProfile ResourceProfile
	Planets         map[string]*Planet
	ZoneType        LayerName
}

type LayerName string

const (
	SectorOne   LayerName = "Sector One"
	SectorTwo   LayerName = "Sector Two"
	SectorThree LayerName = "Sector Three"
	SectorFour  LayerName = "Sector Four"
)

func (w *World) GenerateZones(numZones int) {
	type zonePercentages struct {
		zType      LayerName
		percentage float64
	}

	fmt.Println(w.LayerBoundaries)

	cp := 0
	pcp := 0.0
	for i := 0; i < numZones; i++ {
		currentZone, err := w.GetZoneByIndex(cp)
		if err != nil {
			// TODO: Handle error correctly
			continue
		}
		if float64(numZones)*(pcp+currentZone.MapPercentage) < float64(i) {
			pcp += currentZone.MapPercentage
			cp++
			currentZone, err = w.GetZoneByIndex(cp)
			if err != nil {
				continue
			}
		}

		// Calculate location
		innerRadius := 0.0

		if cp != 0 {
			innerRadius = w.LayerBoundaries[cp-1]
		}
		outerRadius := w.LayerBoundaries[cp]

		radius := innerRadius + (outerRadius-innerRadius)*w.RandomNumber.Float64()
		angle := 2 * math.Pi * w.RandomNumber.Float64()

		x := w.Size/2 + radius*math.Cos(angle)
		y := w.Size/2 + radius*math.Sin(angle)

		dangerLevel := w.randomInt(currentZone.LowerDanger, currentZone.HigherDanger)
		zone := Zone{
			Name:            fmt.Sprintf("Zone-%d", i+1),
			CentralPoint:    Coordinates{x, y},
			DangerRange:     [2]int{dangerLevel, dangerLevel + 10},
			ResourceProfile: GenerateResourceProfile(),
			ZoneType:        LayerName(currentZone.Name),
		}

		planetsAmount := w.randomInt(currentZone.LowerPlanetsAmount, currentZone.HigherPlanetsAmount)

		zone.Planets = w.GeneratePlanetsInZone(planetsAmount, zone, w.AllZoneTypes[zone.ZoneType])

		w.Zones[zone.Name] = &zone
	}
}
