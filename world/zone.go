package world

import "fmt"

type ZoneType struct {
	Name                 ZoneName
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
	ZoneType        ZoneName
}

type ZoneName string

const (
	GuardianSectors ZoneName = "Guardian Sectors"
	TradeLanes      ZoneName = "Trade Lanes"
	OutlawQuadrants ZoneName = "Outlaw Quadrants"
	DarkMatterZones ZoneName = "Dark Matter Zones"
)

func (w *World) GenerateZones(numZones int) {
	type zonePercentages struct {
		zType      ZoneName
		percentage float64
	}
	zp := []zonePercentages{
		{
			zType:      GuardianSectors,
			percentage: 0.30,
		},
		{
			zType:      TradeLanes,
			percentage: 0.30,
		},
		{
			zType:      OutlawQuadrants,
			percentage: 0.25,
		},
		{
			zType:      DarkMatterZones,
			percentage: 0.15,
		},
	}

	cp := 0
	pcp := 0.0

	for i := 0; i < numZones; i++ {

		fmt.Println()

		if float64(numZones)*(pcp+zp[cp].percentage) < float64(i) {
			pcp += zp[cp].percentage
			cp++
		}

		currentZoneType := w.AllZoneTypes[zp[cp].zType]

		dangerLevel := w.randomInt(currentZoneType.LowerDanger, currentZoneType.HigherDanger)
		zone := Zone{
			Name:            fmt.Sprintf("Zone-%d", i+1),
			CentralPoint:    Coordinates{w.RandomNumber.Float64() * 10_000, w.RandomNumber.Float64() * 10_000},
			DangerRange:     [2]int{dangerLevel, dangerLevel + 10},
			ResourceProfile: GenerateResourceProfile(),
			ZoneType:        ZoneName(currentZoneType.Name),
		}

		planetsAmount := w.randomInt(currentZoneType.LowerPlanetsAmount, currentZoneType.HigherPlanetsAmount)

		zone.Planets = w.GeneratePlanetsInZone(planetsAmount, zone, w.AllZoneTypes[zone.ZoneType])

		w.Zones[zone.Name] = &zone
	}

}
