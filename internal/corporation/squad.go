package corporation

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/ship"
	"github.com/luisya22/galactic-exchange/internal/world"
)

type Squad struct {
	Id          uint64
	Ships       *ship.Ship
	CrewMembers []*CrewMember
	Cargo       map[world.Resource]int
	Location    world.Coordinates
	// Officers []Officers   coming soon...
}

func (s Squad) GetHarvestingBonus() int {
	return 1
}

func (c *Corporation) GetSquadReference(squadIndex int) (Squad, error) {
	c.Rw.RLock()
	defer c.Rw.RUnlock()

	if squadIndex >= len(c.Squads) {
		return Squad{}, fmt.Errorf("error: squad not found %v", squadIndex)
	}

	return *c.Squads[squadIndex], nil
}

func (c *Corporation) GetSquad(squadIndex int) (gamecomm.Squad, error) {
	c.Rw.RLock()
	defer c.Rw.RUnlock()

	if squadIndex < 0 {
		return gamecomm.Squad{}, fmt.Errorf("error: invalid squadId: %v", squadIndex)
	}
	if squadIndex >= len(c.Squads) {
		return gamecomm.Squad{}, fmt.Errorf("error: squad not found %v", squadIndex)
	}

	squad := *c.Squads[squadIndex]

	return squad.copy(), nil
}

func (c *Corporation) AddResourceToSquad(squadIndex int, resource world.Resource, amount int) (int, error) {
	var squad *Squad
	c.Rw.Lock()
	defer c.Rw.Unlock()

	if amount < 0 {
		return 0, fmt.Errorf("error: amoutn should be greater than zero")
	}

	if len(c.Squads) <= squadIndex {
		return 0, fmt.Errorf("error: squad not found %v", squadIndex)
	}

	squad = c.Squads[squadIndex]

	squad.Cargo[resource] += amount

	return squad.Cargo[resource], nil
}

func (c *Corporation) RemoveResourcesFromSquad(squadIndex int, resource world.Resource, amount int) (int, error) {
	var squad *Squad
	c.Rw.Lock()
	defer c.Rw.Unlock()

	if amount < 0 {
		return 0, fmt.Errorf("error: amount should be greater than zero")
	}

	if squadIndex >= len(c.Squads) {
		return 0, fmt.Errorf("error: squad not found %v", squadIndex)
	}

	squad = c.Squads[squadIndex]

	if squad.Cargo[resource] < amount {
		return 0, fmt.Errorf("error: squad doesn't have enough amount of resource %v", resource)
	}

	squad.Cargo[resource] -= amount

	return squad.Cargo[resource], nil
}

func (c *Corporation) RemoveAllResourcesFromSquad(squadIndex int, resource world.Resource) (int, error) {
	var squad *Squad

	c.Rw.Lock()
	defer c.Rw.Unlock()

	if squadIndex >= len(c.Squads) {
		return 0, fmt.Errorf("error: squad not found %v", squadIndex)
	}

	squad = c.Squads[squadIndex]

	squad.Cargo[resource] = 0

	return squad.Cargo[resource], nil

}

func (s *Squad) copy() gamecomm.Squad {

	crew := []gamecomm.CrewMember{}
	for _, cm := range s.CrewMembers {
		crew = append(crew, cm.Copy())
	}

	cargo := make(map[string]int)
	for r, c := range s.Cargo {
		cargo[string(r)] = c
	}

	coordinates := gamecomm.Coordinates{X: s.Location.X, Y: s.Location.Y}

	return gamecomm.Squad{
		Id:          s.Id,
		Ships:       s.Ships.Copy(),
		CrewMembers: crew,
		Cargo:       cargo,
		Location:    coordinates,
	}
}
