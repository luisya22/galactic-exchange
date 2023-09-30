package game

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/luisya22/galactic-exchange/corporation"
	"github.com/luisya22/galactic-exchange/world"
)

type Game struct {
	World        *world.World
	Corporations *corporation.CorpGroup
	PlayerState  *PlayerState
}

type Ship struct {
	ID              uint64
	Name            string
	Location        world.Coordinates
	CargoCapacity   float64
	StoredResources map[string]float64
	Crew            []*corporation.CrewMember
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
	Corporation *corporation.Corporation
	Name        string
}

func New() *Game {
	w := world.New()

	playerState := newPlayer()
	corporations := corporation.NewCorpGroup()

	corporations.Corporations[1] = playerState.Corporation

	return &Game{
		World:        w,
		PlayerState:  playerState,
		Corporations: corporations,
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func Start() {
	reader := bufio.NewReader(os.Stdin)

	game := New()

	for k, p := range game.World.Planets {
		fmt.Printf("Planets: %v -> %v -> %v -> %v\n", len(game.World.Planets), k, strconv.Quote(p.Name), game.World.Planets["Zone-1-Planet-1"])
		break
	}

	fmt.Printf("%v\n", game.PlayerState.Corporation.Bases[0].StoredResources[world.Iron])

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
	squadId := command[2]

	return g.HarvestPlanet(planetId, 1, squadId)
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

	playerCorporation := &corporation.Corporation{
		ID:          1,
		Name:        "Player One Corporation",
		Reputation:  0,
		Credits:     10_000,
		Bases:       playerBases,
		CrewMembers: crewMembers,
		IsPlayer:    true,
	}

	return &PlayerState{
		Corporation: playerCorporation,
		Name:        "Player One",
	}

}
