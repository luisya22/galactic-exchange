package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/luisya22/galactic-exchange/world"
)

type Game struct {
	World       *world.World
	PlayerState *PlayerState
}

// TODO: Maybe move to player package later
type CorporationState struct {
	ID                              uint64
	Name                            string
	Reputation                      int
	Credits                         float64
	Bases                           []*Base
	CrewMembers                     []*CrewMember
	IsPlayer                        bool
	ReputationWithOtherCorporations map[string]int
}

type Base struct {
	ID                 uint64
	Name               string
	Location           world.Coordinates
	ResourceProduction map[world.Resource]int
	StorageCapacity    float64
	StoredResources    map[world.Resource]int
}

type Ship struct {
	ID              uint64
	Name            string
	Location        world.Coordinates
	CargoCapacity   float64
	StoredResources map[string]float64
	Crew            []*CrewMember
}

type CrewMember struct {
	ID         uint64
	Name       string
	Species    string
	Skills     map[string]int
	AssignedTo uint64
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
	Corporation *CorporationState
	Name        string
}

func New() *Game {
	w := world.New()

	playerState := newPlayer()

	return &Game{
		World:       w,
		PlayerState: playerState,
	}
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

	// <command> args...
	for {
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
				fmt.Printf("Wrong command: the sell command is 'sell <amount> <item> <planetId>")
				continue
			}

			err := game.sellResource(command)
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
	var amount int
	var item world.ResourceInfo
	var planet *world.Planet
	var playerResourceAmount int
	var ok bool

	amount, err := strconv.Atoi(command[1])
	if err != nil {
		return fmt.Errorf("Wrong command: the second arg should be a number on the sell command")
	}

	if item, ok = g.World.AllResources[world.Resource(command[2])]; !ok {
		return fmt.Errorf("Item not found")
	}

	if planet, ok = g.World.Planets[command[3]]; !ok {
		return fmt.Errorf("Planet not found: %v -> %v", strconv.Quote(command[3]), planet)
	}

	// Check that the player has enough of an item
	if playerResourceAmount, ok = g.PlayerState.Corporation.Bases[0].StoredResources[item.Name]; !ok {
		return fmt.Errorf("You don't have enough %v: 0", item.Name)
	}

	if playerResourceAmount < amount {
		return fmt.Errorf("You don't have enough %s: %d", item.Name, playerResourceAmount)
	}

	//Save old amounts
	planetAmount := planet.Resources[item.Name]
	// Transfer amounts
	g.PlayerState.Corporation.Bases[0].StoredResources[item.Name] -= amount
	g.World.Planets[planet.Name].Resources[item.Name] += amount
	g.PlayerState.Corporation.Credits += item.BasePrice * float64(amount)
	//Print transfer
	fmt.Printf(
		"Transfer:\nPlanet: %v -> %v\nCorporation: %v -> %v\nSell Price: %v\nNew Player Credits Balance: %v\n",
		planetAmount,
		g.World.Planets[planet.Name].Resources[item.Name],
		playerResourceAmount,
		g.PlayerState.Corporation.Bases[0].StoredResources[item.Name],
		item.BasePrice,
		g.PlayerState.Corporation.Credits,
	)

	return nil
}

func newPlayer() *PlayerState {

	playerBases := []*Base{
		{
			ID:              1,
			Name:            "Player One Base",
			Location:        world.Coordinates{X: 0, Y: 0},
			StorageCapacity: 50_000,
			StoredResources: map[world.Resource]int{world.Iron: 1000},
		},
	}

	crewMembers := []*CrewMember{
		{
			ID:         1,
			Name:       "Galios Trek",
			Species:    "Bertusian",
			AssignedTo: 1,
		},
	}

	playerCorporation := &CorporationState{
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
