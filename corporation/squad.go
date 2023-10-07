package corporation

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/ship"
	"github.com/luisya22/galactic-exchange/world"
)

type Squad struct {
	Ships       *ship.Ship
	CrewMembers []*CrewMember
	Cargo       map[world.Resource]int
	// Officers []Officers   coming soon...
}

func (c *Corporation) GetSquad(squadIndex int) (Squad, error) {
	c.Rw.RLock()
	defer c.Rw.RUnlock()

	if len(c.Squads) > squadIndex {
		return *c.Squads[squadIndex], fmt.Errorf("error: squad not found %v", squadIndex)
	}

	return Squad{}, nil
}

func (c *Corporation) AddResourceToSquad(squadIndex int, resource world.Resource, amount int) (int, error) {
	var squad *Squad
	c.Rw.Lock()
	defer c.Rw.Unlock()

	if len(c.Squads) > squadIndex {
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

	if len(c.Squads) > squadIndex {
		return 0, fmt.Errorf("error: squad not found %v", squadIndex)
	}

	squad = c.Squads[squadIndex]

	if squad.Cargo[resource] == 0 {
		return 0, fmt.Errorf("error: squad doesn't have enough amount of resource %v", resource)
	}

	squad.Cargo[resource] -= amount

	return squad.Cargo[resource], nil
}
