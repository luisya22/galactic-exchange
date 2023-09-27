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

func (w *World) GenerateZones(mapSize float64, numZones int) {
	type zonePercentages struct {
		zType      LayerName
		percentage float64
	}
	// Define layer percentages
	zp := []zonePercentages{
		{
			zType:      SectorOne,
			percentage: 0.15,
		},
		{
			zType:      SectorTwo,
			percentage: 0.20,
		},
		{
			zType:      SectorThree,
			percentage: 0.20,
		},
		{
			zType:      SectorFour,
			percentage: 0.45,
		},
	}

	layerBoundaries := []float64{}

	// Calculate layer boundaries
	currentBoundary := 0.0
	for _, zonePercentage := range zp {
		currentBoundary += mapSize / 2 * zonePercentage.percentage

		layerBoundaries = append(layerBoundaries, currentBoundary)
	}

	cp := 0
	pcp := 0.0
	for i := 0; i < numZones; i++ {
		if float64(numZones)*(pcp+zp[cp].percentage) < float64(i) {
			pcp += zp[cp].percentage
			cp++
		}

		// Calculate location
		innerRadius := 0.0

		if cp != 0 {
			innerRadius = layerBoundaries[cp-1]
		}
		outerRadius := layerBoundaries[cp]

		// if cp+1 != len(layerBoundaries) {
		// 	outerRadius = layerBoundaries[cp+1]
		// }

		radius := innerRadius + (outerRadius-innerRadius)*w.RandomNumber.Float64()
		angle := 2 * math.Pi * w.RandomNumber.Float64()

		x := mapSize/2 + radius*math.Cos(angle)
		y := mapSize/2 + radius*math.Sin(angle)

		currentZoneType := w.AllZoneTypes[zp[cp].zType]

		dangerLevel := w.randomInt(currentZoneType.LowerDanger, currentZoneType.HigherDanger)
		zone := Zone{
			Name:            fmt.Sprintf("Zone-%d", i+1),
			CentralPoint:    Coordinates{x, y},
			DangerRange:     [2]int{dangerLevel, dangerLevel + 10},
			ResourceProfile: GenerateResourceProfile(),
			ZoneType:        LayerName(currentZoneType.Name),
		}

		planetsAmount := w.randomInt(currentZoneType.LowerPlanetsAmount, currentZoneType.HigherPlanetsAmount)

		zone.Planets = w.GeneratePlanetsInZone(planetsAmount, zone, w.AllZoneTypes[zone.ZoneType])

		w.Zones[zone.Name] = &zone
	}
}

// func (w *World) GenerateZones(numZones int) {
// 	type zonePercentages struct {
// 		zType      ZoneName
// 		percentage float64
// 	}
// 	zp := []zonePercentages{
// 		{
// 			zType:      GuardianSectors,
// 			percentage: 0.15,
// 		},
// 		{
// 			zType:      TradeLanes,
// 			percentage: 0.20,
// 		},
// 		{
// 			zType:      OutlawQuadrants,
// 			percentage: 0.20,
// 		},
// 		{
// 			zType:      DarkMatterZones,
// 			percentage: 0.55,
// 		},
// 	}
//
// 	cp := 0
// 	pcp := 0.0
//
// 	for i := 0; i < numZones; i++ {
// 		if float64(numZones)*(pcp+zp[cp].percentage) < float64(i) {
// 			pcp += zp[cp].percentage
// 			cp++
// 		}
//
// 		currentZoneType := w.AllZoneTypes[zp[cp].zType]
//
// 		dangerLevel := w.randomInt(currentZoneType.LowerDanger, currentZoneType.HigherDanger)
// 		zone := Zone{
// 			Name:            fmt.Sprintf("Zone-%d", i+1),
// 			CentralPoint:    Coordinates{w.RandomNumber.Float64() * 10_000, w.RandomNumber.Float64() * 10_000},
// 			DangerRange:     [2]int{dangerLevel, dangerLevel + 10},
// 			ResourceProfile: GenerateResourceProfile(),
// 			ZoneType:        ZoneName(currentZoneType.Name),
// 		}
//
// 		planetsAmount := w.randomInt(currentZoneType.LowerPlanetsAmount, currentZoneType.HigherPlanetsAmount)
//
// 		zone.Planets = w.GeneratePlanetsInZone(planetsAmount, zone, w.AllZoneTypes[zone.ZoneType])
//
// 		w.Zones[zone.Name] = &zone
// 	}
//
// }
