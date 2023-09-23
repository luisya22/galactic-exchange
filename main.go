package main

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/world"
)

func main() {

	world := world.NewWorld()

	for _, z := range world.Zones {
		fmt.Printf(
			"------\nZONE: %s\nZone Type: %s\nCentral Point: %v\nDanger Range: %v\nResource Profile: %v\n\nPlanets:\n",
			z.Name,
			z.ZoneType,
			z.CentralPoint,
			z.DangerRange,
			z.ResourceProfile,
		)

		for _, p := range z.Planets {
			fmt.Printf(
				"%v:\nLocation: %v\nResources:%v\nPopulation:%v\nDangerLevel:%v\n-------\n",
				p.Name,
				p.Location,
				p.Resources,
				p.Population,
				p.DangerLevel,
			)
		}

		fmt.Println("---------")

	}

	PlotZonesASCII(*world)
}

const gridWidth, gridHeight = 200, 200
const universeSize = 10_000.0
const cellSize = universeSize / float64(gridWidth) // How many universe units are in each grid cell

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func PlotZonesASCII(world world.World) {
	// Initialize grid with empty space
	grid := make([][]rune, gridHeight)
	for i := range grid {
		grid[i] = make([]rune, gridWidth)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	for _, z := range world.Zones {
		// Plotting the central point of the zone
		zX, zY := int(z.CentralPoint.X/cellSize), int(z.CentralPoint.Y/cellSize)
		zX = clamp(zX, 0, gridWidth-1)
		zY = clamp(zY, 0, gridHeight-1)
		grid[zY][zX] = 'O'

		// Plotting the planets within the zone
		for _, planet := range z.Planets {
			pX, pY := int(planet.Location.X/cellSize), int(planet.Location.Y/cellSize)
			pX = clamp(pX, 0, gridWidth-1)
			pY = clamp(pY, 0, gridHeight-1)
			if grid[pY][pX] != 'O' { // We avoid overriding a zone center
				grid[pY][pX] = '0'
			}
		}
	}

	// Print the grid
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}
