package corporation_test

import (
	"testing"

	"github.com/luisya22/galactic-exchange/internal/corporation"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/ship"
	"github.com/luisya22/galactic-exchange/internal/world"
)

const (
	initialCorporationCredits = 10_000
	corporationName           = "Test Corporation"
	corporationID             = 1
	baseName                  = "Test Base"
	testSquadId               = 14
	initialIronQuantity       = 1000
)

func createTestCorpGroup(t *testing.T, gameChannels *gamecomm.GameChannels) *corporation.CorpGroup {

	t.Helper()

	corporation1 := createTestCorporation(t)

	corporations := make(map[uint64]*corporation.Corporation)
	corporations[corporation1.ID] = corporation1

	return &corporation.CorpGroup{
		Corporations: corporations,
		Workers:      10,
		CorpChan:     gameChannels.CorpChannel,
	}
}

func createTestCorporation(t *testing.T) *corporation.Corporation {

	t.Helper()

	playerBases := []*corporation.Base{
		{
			ID:              1,
			Name:            baseName,
			Location:        world.Coordinates{X: 0, Y: 0},
			StorageCapacity: 50_000,
			StoredResources: map[string]int{"iron": initialIronQuantity},
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
	}

	squads := []*corporation.Squad{
		{
			Id:          testSquadId,
			Ships:       ship,
			CrewMembers: []*corporation.CrewMember{crewMembers[0]},
			Cargo:       make(map[string]int),
		},
	}

	squads[0].Cargo["iron"] = initialIronQuantity

	playerCorporation := &corporation.Corporation{
		ID:          corporationID,
		Name:        corporationName,
		Reputation:  0,
		Credits:     initialCorporationCredits,
		Bases:       playerBases,
		CrewMembers: crewMembers,
		IsPlayer:    true,
		Squads:      squads,
	}

	return playerCorporation
}
