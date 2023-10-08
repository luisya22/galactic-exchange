package world

import "math"

type Coordinates struct {
	X float64
	Y float64
}

func Distance(p1, p2 Coordinates) float64 {
	dx := math.Pow(p2.X-p1.X, 2)
	dy := math.Pow(p2.Y-p1.Y, 2)

	return math.Sqrt(dx + dy)
}
