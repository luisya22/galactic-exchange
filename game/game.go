package game

import "github.com/luisya22/galactic-exchange/world"

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
	ResourceProduction map[world.Resource]float64
	StorageCapacity    float64
	StoredResources    map[world.Resource]float64
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

func newPlayer() *PlayerState {

	playerBases := []*Base{
		{
			ID:              1,
			Name:            "Player One Base",
			Location:        world.Coordinates{X: 0, Y: 0},
			StorageCapacity: 50_000,
			StoredResources: map[world.Resource]float64{world.Iron: 1000},
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
