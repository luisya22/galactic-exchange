package main

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/game"
	"github.com/luisya22/galactic-exchange/world"
)

func main() {
	game.Start()
}

const gridWidth, gridHeight = 500, 500
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

func PlotZonesASCII(w world.World) {
	// Initialize grid with empty space
	grid := make([][]string, gridHeight)
	for i := range grid {
		grid[i] = make([]string, gridWidth)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	for _, z := range w.Zones {
		// Plotting the central point of the zone
		zX, zY := int(z.CentralPoint.X/cellSize), int(z.CentralPoint.Y/cellSize)
		zX = clamp(zX, 0, gridWidth-1)
		zY = clamp(zY, 0, gridHeight-1)

		switch z.ZoneType {
		case world.SectorOne:
			grid[zY][zX] = "\033[34mG\033[0m"
		case world.SectorTwo:
			grid[zY][zX] = "\033[32mT\033[0m"
		case world.SectorThree:
			grid[zY][zX] = "\033[33mO\033[0m"
		case world.SectorFour:
			grid[zY][zX] = "\033[31mD\033[0m"

		}

		// Plotting the planets within the zone
		for _, planet := range z.Planets {
			pX, pY := int(planet.Location.X/cellSize), int(planet.Location.Y/cellSize)
			pX = clamp(pX, 0, gridWidth-1)
			pY = clamp(pY, 0, gridHeight-1)

			switch z.ZoneType {
			case world.SectorOne:
				grid[pY][pX] = "\033[34mG\033[0m"
			case world.SectorTwo:
				grid[pY][pX] = "\033[32mT\033[0m"
			case world.SectorThree:
				grid[pY][pX] = "\033[33mO\033[0m"
			case world.SectorFour:
				grid[pY][pX] = "\033[31mD\033[0m"

			}
		}
	}

	grid[int((universeSize/2)/cellSize)][int((universeSize/2)/cellSize)] = "X"

	// Print the grid
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%s", cell)
		}
		fmt.Println()
	}
}
