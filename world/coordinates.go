package world

import "math"

type Coordinates struct {
	X float64
	Y float64
}

func Distance(p1, p2 Planet) float64 {
	dx := math.Pow(p2.Location.X-p1.Location.X, 2)
	dy := math.Pow(p2.Location.Y-p1.Location.Y, 2)

	return math.Sqrt(dx + dy)
}
