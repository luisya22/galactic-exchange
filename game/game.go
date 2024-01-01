package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/luisya22/galactic-exchange/corporation"
	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/mission"
	"github.com/luisya22/galactic-exchange/ship"
	"github.com/luisya22/galactic-exchange/world"
)

// TODO: Gracefully shutdown game
type Game struct {
	World            *world.World
	Corporations     *corporation.CorpGroup
	PlayerState      *PlayerState
	MissionScheduler *mission.MissionScheduler
	gameChannels     *gamecomm.GameChannels
	gameClock        *gameclock.GameClock
}

/*
Possible State fields:

	Achievements
	GameProgress
	ActiveQuests
	Skills
	Inventory
	Ship and building Blueprints
	LastSavedTime
	Etc...
*/
type PlayerState struct {
	Corporation      *corporation.Corporation
	Name             string
	NotificationChan chan string
}

func (ps *PlayerState) listenNotifications() {
	for message := range ps.NotificationChan {
		log.Println(message)
	}
}

func New() *Game {
	gameChannels := &gamecomm.GameChannels{
		WorldChannel:   make(chan gamecomm.WorldCommand, 100),
		CorpChannel:    make(chan gamecomm.CorpCommand, 100),
		MissionChannel: make(chan gamecomm.MissionCommand, 100),
	}

	gc := gameclock.NewGameClock(0, 1)

	w := world.New(gameChannels)

	playerState := newPlayer()
	corporations := corporation.NewCorpGroup(gameChannels)

	corporations.Corporations[1] = playerState.Corporation

	missionScheduler := mission.NewMissionScheduler(gameChannels, gc)

	return &Game{
		World:            w,
		PlayerState:      playerState,
		Corporations:     corporations,
		MissionScheduler: missionScheduler,
		gameChannels:     gameChannels,
		gameClock:        gc,
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func Start() error {
	reader := bufio.NewReader(os.Stdin)

	game := New()

	PlotZonesASCII(*game.World)

	go game.MissionScheduler.Run()
	go game.Corporations.Run()
	go game.PlayerState.listenNotifications()
	go game.gameClock.StartTime()

	for k, p := range game.World.Planets {
		fmt.Printf("Planets: %v -> %v -> %v -> %v\n", len(game.World.Planets), k, strconv.Quote(p.Name), game.World.Planets["Zone-1-Planet-1"])
		break
	}

	fmt.Printf("%v\n", game.PlayerState.Corporation.Bases[0].StoredResources[world.Iron])
	fmt.Println(game.PlayerState.Corporation.Squads[0])

	for _, i := range game.World.AllResources {
		fmt.Printf("%v -> %v\n", i.Name, i.BasePrice)
	}

	var memStats runtime.MemStats

	// Collect memory stats

	// <command> args...
	for {
		runtime.ReadMemStats(&memStats)

		// Print memory stats
		fmt.Printf("Alloc = %v MiB", bToMb(memStats.Alloc))
		fmt.Printf("\tTotalAlloc = %v MiB", bToMb(memStats.TotalAlloc))
		fmt.Printf("\tSys = %v MiB", bToMb(memStats.Sys))
		fmt.Printf("\tNumGC = %v\n", memStats.NumGC)
		fmt.Printf("\tGoRoutines = %v\n", runtime.NumGoroutine())

		fmt.Println()
		fmt.Print("Enter command: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Problem reading input: %v", err.Error())
			continue
		}

		str := strings.ReplaceAll(input, "\n", "")
		command := strings.Split(str, " ")

		switch command[0] {

		case "sell":
			if len(command) != 4 {
				fmt.Printf("Wrong command: the sell command is 'sell <amount> <item> <planetId>'")
				continue
			}

			err := game.sellResource(command)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "harvest":
			if len(command) != 3 {
				fmt.Printf("Wrong command: the harvest command is 'harvets <planet> <squad>'")
				continue
			}

			err := game.harvestPlanet(command)
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			fmt.Printf("Wrong command %v\n", command)
		}
	}
}

// sell <number> <item> <planet>
func (g *Game) sellResource(command []string) error {

	amount, err := strconv.Atoi(command[1])
	if err != nil {
		return fmt.Errorf("%v needs to be an integer", command[1])
	}

	itemName := world.Resource(command[2])
	planetId := command[3]

	return g.SellResource(amount, itemName, planetId, 1)

}

// harvest <planet> <squad>
func (g *Game) harvestPlanet(command []string) error {
	planetId := command[1]
	squadId, err := strconv.Atoi(command[2])
	if err != nil {
		return fmt.Errorf("%v needs to be an integer", command[2])
	}

	return g.HarvestPlanet(planetId, 1, squadId, g.PlayerState.NotificationChan)
}

func newPlayer() *PlayerState {

	playerBases := []*corporation.Base{
		{
			ID:              1,
			Name:            "Player One Base",
			Location:        world.Coordinates{X: 0, Y: 0},
			StorageCapacity: 50_000,
			StoredResources: map[world.Resource]int{world.Iron: 1000},
		},
	}

	crewMembers := []*corporation.CrewMember{
		{
			ID:         1,
			Name:       "Galios Trek",
			Species:    "Bertusian",
			AssignedTo: 1,
		},
	}

	shipLocation := world.Coordinates{
		X: playerBases[0].Location.X,
		Y: playerBases[0].Location.Y,
	}

	ship := &ship.Ship{
		Name:         "MF",
		Capacity:     10,
		MaxHealth:    1000,
		ActualHealth: 1000,
		MaxCargo:     10_000,
		Location:     shipLocation,
		Speed:        10,
		// Attributes
		// Upgrades
		// StoredResources
	}

	squads := []*corporation.Squad{
		{
			Ships:       ship,
			CrewMembers: []*corporation.CrewMember{crewMembers[0]},
			Cargo:       make(map[world.Resource]int),
		},
	}

	playerCorporation := &corporation.Corporation{
		ID:          1,
		Name:        "Player One Corporation",
		Reputation:  0,
		Credits:     10_000,
		Bases:       playerBases,
		CrewMembers: crewMembers,
		IsPlayer:    true,
		Squads:      squads,
	}

	return &PlayerState{
		Corporation:      playerCorporation,
		Name:             "Player One",
		NotificationChan: make(chan string),
	}

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
