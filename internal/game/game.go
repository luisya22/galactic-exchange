package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/luisya22/galactic-exchange/internal/corporation"
	"github.com/luisya22/galactic-exchange/internal/economy"
	"github.com/luisya22/galactic-exchange/internal/gameclock"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/mission"
	"github.com/luisya22/galactic-exchange/internal/resource"
	"github.com/luisya22/galactic-exchange/internal/ship"
	"github.com/luisya22/galactic-exchange/internal/world"
)

// TODO: Gracefully shutdown game
type Game struct {
	World            *world.World
	Corporations     *corporation.CorpGroup
	PlayerState      *PlayerState
	MissionScheduler *mission.MissionScheduler
	gameChannels     *gamecomm.GameChannels
	gameClock        *gameclock.GameClock
	Resources        map[string]resource.Resource
	Economy          *economy.Economy
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
		EconomyChannel: make(chan gamecomm.EconomyCommand, 100),
	}

	resources := resource.LoadWorldResources()

	gc := gameclock.NewGameClock(0, 1)

	w := world.New(gameChannels, resources, gc)

	playerState := newPlayer()
	corporations := corporation.NewCorpGroup(gameChannels)
	gameEconomy := economy.NewEconomy(*gameChannels, resources, w.GetZoneIds(), gc)

	corporations.Corporations[1] = playerState.Corporation

	missionScheduler := mission.NewMissionScheduler(gameChannels, gc)

	return &Game{
		World:            w,
		PlayerState:      playerState,
		Corporations:     corporations,
		MissionScheduler: missionScheduler,
		gameChannels:     gameChannels,
		gameClock:        gc,
		Resources:        resource.LoadWorldResources(),
		Economy:          gameEconomy,
	}
}

func Start() error {
	reader := bufio.NewReader(os.Stdin)

	game := New()

	go game.MissionScheduler.Run()
	go game.Corporations.Run()
	go game.PlayerState.listenNotifications()
	go game.gameClock.StartTime()
	go game.Economy.Run()

	printTestLog(game)

	// <command> args...
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Problem reading input: %v", err.Error())
			continue
		}

		str := strings.ReplaceAll(input, "\n", "")
		command := strings.Split(str, " ")

		switch command[0] {

		case "sell":
			if len(command) != 5 {
				fmt.Printf("Wrong command: the sell command is 'sell <amount> <item> <planetId> <squad>'")
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

// sell <number> <item> <planet> <squadId>
func (g *Game) sellResource(command []string) error {

	amount, err := strconv.Atoi(command[1])
	if err != nil {
		return fmt.Errorf("%v needs to be an integer", command[1])
	}

	itemName := command[2]
	planetId := command[3]

	squadId, err := strconv.Atoi(command[4])
	if err != nil {
		return fmt.Errorf("%v needs to be an integer", command[4])
	}

	return g.SellResource(planetId, 1, squadId, amount, itemName, g.PlayerState.NotificationChan)
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
			StoredResources: map[string]int{"iron": 1000},
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
			Cargo:       make(map[string]int),
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
