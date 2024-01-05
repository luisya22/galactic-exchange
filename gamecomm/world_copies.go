package gamecomm

import (
	"math"
)

type Planet struct {
	Name           string
	Location       Coordinates
	Resources      map[string]int
	Population     int
	DangerLevel    int
	ResourceDemand map[string]int
	IsHabitable    bool
	IsHarvestable  bool
}

type Coordinates struct {
	X float64
	Y float64
}

func Distance(p1, p2 Coordinates) float64 {
	dx := math.Pow(p2.X-p1.X, 2)
	dy := math.Pow(p2.Y-p1.Y, 2)

	return math.Sqrt(dx + dy)
}
