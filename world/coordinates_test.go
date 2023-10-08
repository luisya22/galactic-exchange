package world_test

import (
	"math"
	"testing"

	"github.com/luisya22/galactic-exchange/world"
)

const epsilon = 1e-6

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func TestDistance(t *testing.T) {

	testMap := []struct {
		name    string
		planet1 *world.Planet
		planet2 *world.Planet
		wants   float64
	}{
		{
			name: "Positive distances - planet one smaller",
			planet1: &world.Planet{
				Name:     "planet1",
				Location: world.Coordinates{25, 25},
			},
			planet2: &world.Planet{
				Name:     "planet2",
				Location: world.Coordinates{100, 100},
			},
			wants: 106.066017,
		},
		{
			name: "Positive distances - planet one bigger",
			planet1: &world.Planet{
				Name:     "planet1",
				Location: world.Coordinates{253, 554},
			},
			planet2: &world.Planet{
				Name:     "planet2",
				Location: world.Coordinates{34, 115},
			},
			wants: 490.593518,
		},
	}

	for _, tt := range testMap {
		t.Run(tt.name, func(t *testing.T) {

			dis := world.Distance(*tt.planet1, *tt.planet2)

			if !almostEqual(dis, tt.wants) {
				t.Errorf("error: wants %v; got: %v; diff: %v", tt.wants, dis, math.Abs(dis-tt.wants))
			}
		})
	}
}
