package corporation

import (
	"fmt"
	"sync"

	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/maputils"
)

type CorpGroup struct {
	Corporations map[uint64]*Corporation
	RW           sync.RWMutex
	Workers      int
	CorpChan     chan gamecomm.CorpCommand
}

type Corporation struct {
	ID                              uint64
	Name                            string
	Reputation                      int
	Credits                         float64
	Bases                           []*Base
	CrewMembers                     []*CrewMember
	Squads                          []*Squad
	IsPlayer                        bool
	ReputationWithOtherCorporations map[string]int
	Rw                              sync.RWMutex
}

func NewCorpGroup(gameChannels *gamecomm.GameChannels) *CorpGroup {
	return &CorpGroup{
		Corporations: make(map[uint64]*Corporation, 50),
		Workers:      100,
		CorpChan:     gameChannels.CorpChannel,
	}
}

func (cg *CorpGroup) Run() {
	cg.Listen()
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

func (c *CorpGroup) findCorporation(corporationId uint64) (gamecomm.Corporation, error) {
	var corporation *Corporation
	var ok bool

	c.RW.RLock()
	defer c.RW.RUnlock()

	if corporation, ok = c.Corporations[corporationId]; !ok {
		return gamecomm.Corporation{}, fmt.Errorf("Corporation not found: %v", corporationId)
	}

	corpCopy := gamecomm.Corporation{
		ID:         corporation.ID,
		Name:       corporation.Name,
		Reputation: corporation.Reputation,
		Credits:    corporation.Credits,
	}

	return corpCopy, nil

}

func (c *Corporation) copy() Corporation {
	basesCopy := make([]*Base, len(c.Bases))
	copy(c.Bases, c.Bases)

	crewMembersCopy := make([]*CrewMember, len(c.CrewMembers))
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

	if amount < 0 {
		return 0, fmt.Errorf("error: amount should be greater than zero")
	}
	corporation, err := c.findCorporationReference(corporationId)
	if err != nil {
		return 0, err
	}

	corporation.Rw.Lock()
	defer corporation.Rw.Unlock()

	corporation.Credits += amount

	return corporation.Credits, nil
}

func (c *CorpGroup) RemoveCredits(corporationId uint64, amount float64) (float64, error) {

	if amount < 0 {
		return 0, fmt.Errorf("error: amount should be greater than zero")
	}

	corporation, err := c.findCorporationReference(corporationId)
	if err != nil {
		return 0, err
	}

	corporation.Rw.Lock()
	defer corporation.Rw.Unlock()

	if amount > corporation.Credits {
		return 0, fmt.Errorf("error: not enough credits")
	}

	corporation.Credits -= amount

	return corporation.Credits, nil
}

func (c *CorpGroup) RemoveResources(corporationId uint64, resource string, amount int) (int, error) {

	if amount < 0 {
		return 0, fmt.Errorf("error: amount should be greater than zero")
	}
	corporation, err := c.findCorporationReference(corporationId)
	if err != nil {
		return 0, err
	}

	corporation.Rw.Lock()
	defer corporation.Rw.Unlock()

	// TODO: Select the correct base
	if resourcesAmount, ok := corporation.Bases[0].StoredResources[resource]; !ok || resourcesAmount < amount {
		return 0, fmt.Errorf("error: not enough resources on base")
	}

	corporation.Bases[0].StoredResources[resource] -= amount

	return corporation.Bases[0].StoredResources[resource], nil
}

func (c *CorpGroup) AddResources(corporationId uint64, resource string, amount int) (int, error) {

	if amount < 0 {
		return 0, fmt.Errorf("error: amount should be greater than zero")
	}
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

	corporation.Bases[0].StoredResources[resource] += amount

	return corporation.Bases[0].StoredResources[resource], nil
}
