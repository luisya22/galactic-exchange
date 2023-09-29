package corporation

import (
	"fmt"
	"sync"

	"github.com/luisya22/galactic-exchange/base"
	"github.com/luisya22/galactic-exchange/crew"
	"github.com/luisya22/galactic-exchange/internal/maputils"
	"github.com/luisya22/galactic-exchange/world"
)

type CorpGroup struct {
	Corporations map[uint64]*Corporation
	RW           sync.RWMutex
}

type Corporation struct {
	ID                              uint64
	Name                            string
	Reputation                      int
	Credits                         float64
	Bases                           []*base.Base
	CrewMembers                     []*crew.CrewMember
	IsPlayer                        bool
	ReputationWithOtherCorporations map[string]int
	Rw                              sync.RWMutex
}

func NewCorpGroup() *CorpGroup {
	return &CorpGroup{
		Corporations: make(map[uint64]*Corporation, 50),
	}
}

func (c *CorpGroup) FindCorporation(corporationId uint64) (Corporation, error) {
	var corporation *Corporation
	var ok bool

	c.RW.RLock()
	defer c.RW.RUnlock()

	if corporation, ok = c.Corporations[corporationId]; !ok {
		return Corporation{}, fmt.Errorf("Corporation not found: %v", corporationId)
	}

	return corporation.copy(), nil
}

func (c *CorpGroup) findCorporationReference(corporationId uint64) (*Corporation, error) {
	var corporation *Corporation
	var ok bool

	c.RW.RLock()
	defer c.RW.RUnlock()

	if corporation, ok = c.Corporations[corporationId]; !ok {
		return nil, fmt.Errorf("Corporation not found: %v", corporationId)
	}

	return corporation, nil
}

func (c *Corporation) copy() Corporation {
	basesCopy := make([]*base.Base, len(c.Bases))
	copy(c.Bases, c.Bases)

	crewMembersCopy := make([]*crew.CrewMember, len(c.CrewMembers))
	copy(c.CrewMembers, c.CrewMembers)

	return Corporation{
		ID:                              c.ID,
		Name:                            c.Name,
		Reputation:                      c.Reputation,
		Credits:                         c.Credits,
		Bases:                           basesCopy,
		CrewMembers:                     crewMembersCopy,
		IsPlayer:                        c.IsPlayer,
		ReputationWithOtherCorporations: maputils.CopyMap(c.ReputationWithOtherCorporations),
	}
}

func (c *CorpGroup) AddCredits(corporationId uint64, amount float64) (float64, error) {
	corporation, err := c.findCorporationReference(corporationId)
	if err != nil {
		return 0, err
	}

	corporation.Rw.Lock()
	defer corporation.Rw.Unlock()

	corporation.Credits += amount

	return corporation.Credits, nil
}

func (c *CorpGroup) RemoveResources(corporationId uint64, resource world.Resource, amount int) (int, error) {
	corporation, err := c.findCorporationReference(corporationId)
	if err != nil {
		return 0, err
	}

	corporation.Rw.Lock()
	defer corporation.Rw.Unlock()

	// TODO: Select the correct base
	if resourcesAmount, ok := corporation.Bases[0].StoredResources[resource]; !ok || resourcesAmount < amount {
		return 0, fmt.Errorf("Not enough resources on base")
	}

	corporation.Bases[0].StoredResources[resource] -= amount

	return corporation.Bases[0].StoredResources[resource], nil
}

func (c *CorpGroup) AddResources(corporationId uint64, resource world.Resource, amount int) (int, error) {
	corporation, err := c.findCorporationReference(corporationId)
	if err != nil {
		return 0, err
	}

	corporation.Rw.Lock()
	defer corporation.Rw.Unlock()

	// TODO: Select the correct base
	if _, ok := corporation.Bases[0].StoredResources[resource]; !ok {
		corporation.Bases[0].StoredResources[resource] = 0
	}

	corporation.Bases[0].StoredResources[resource] -= amount

	return corporation.Bases[0].StoredResources[resource], nil
}
