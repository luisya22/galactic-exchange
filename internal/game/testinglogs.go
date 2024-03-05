package game

import (
	"fmt"
	"strconv"

	"github.com/luisya22/galactic-exchange/internal/world"
)

// func bToMb(b uint64) uint64 {
// 	return b / 1024 / 1024
// }

func printTestLog(g *Game) {
	for k, p := range g.World.Planets {
		fmt.Printf("Planets: %v -> %v -> %v -> %v\n", len(g.World.Planets), k, strconv.Quote(p.Name), g.World.Planets["Zone-1-Planet-1"])
		break
	}

	fmt.Printf("%v\n", g.PlayerState.Corporation.Bases[0].StoredResources["iron"])
	fmt.Println(g.PlayerState.Corporation.Squads[0])

	for _, i := range g.World.AllResources {
		fmt.Printf("%v -> %v\n", i.Name, i.BasePrice)

	}
}

// var memStats runtime.MemStats
//
// func printBenchmarkData(memStats runtime.MemStats) {
//
// 	runtime.ReadMemStats(&memStats)
//
// 	// Print memory stats
// 	fmt.Printf("Alloc = %v MiB", bToMb(memStats.Alloc))
// 	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(memStats.TotalAlloc))
// 	fmt.Printf("\tSys = %v MiB", bToMb(memStats.Sys))
// 	fmt.Printf("\tNumGC = %v\n", memStats.NumGC)
// 	fmt.Printf("\tGoRoutines = %v\n", runtime.NumGoroutine())
//
// 	fmt.Println()
// }

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

func PlotZonesASCII(w *world.World) {
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
		default:
			grid[zY][zX] = "C"
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
