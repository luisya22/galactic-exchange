package economy

import (
	"sync"

	"github.com/luisya22/galactic-exchange/internal/gameclock"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

// Communicate with World to check planet demands
// Communicate with World to check planet financial budget
// Communicate with World to get Resource Scarcity
// Communicate with World to check planet events or emergencies
// Calculate current market trends

type Economy struct {
	transactions                   []transaction
	zoneTransactions               map[string][]int
	planetTransactions             map[string][]int
	corporationPlanetTradeRelation map[uint64]int
	contracts                      []contract
	corporationContracts           map[uint64][]int
	gameChannels                   gamecomm.GameChannels
	Workers                        int
	rw                             sync.Mutex
}

type transaction struct {
	planetId      string
	corporationId uint64
	resource      string
	credits       float64
	time          gameclock.GameTimeDuration
}

type contract struct {
	corporationId uint64
	planetId      string
	resource      string
	price         float64
	interval      gameclock.GameTimeDuration
	endTime       gameclock.GameTimeDuration
}

type itemAndPrice struct {
	resource string
	price    int
}

func NewEconomy(gameChannels gamecomm.GameChannels) *Economy {
	return &Economy{
		transactions:                   []transaction{},
		zoneTransactions:               make(map[string][]int),
		corporationPlanetTradeRelation: make(map[uint64]int),
		contracts:                      []contract{},
		corporationContracts:           make(map[uint64][]int),
		gameChannels:                   gameChannels,
	}
}

// Decide wether a Planet accepts or decline a trade
func (e *Economy) acceptTrade(planetId string, corporationId uint64, itemPrices []itemAndPrice) bool {
	return true
}

// Save Transactions per zone and globally
func (e *Economy) addTransaction(zoneId string, planetId string, corporationId uint64, resource string, credits float64, t gameclock.GameTimeDuration) error {
	e.rw.Lock()
	defer e.rw.Unlock()

	tran := transaction{
		planetId:      planetId,
		corporationId: corporationId,
		resource:      resource,
		credits:       credits,
		time:          t,
	}

	index := len(e.transactions)

	// Transaction
	e.transactions = append(e.transactions, tran)

	// Zone Transaction
	_, ok := e.zoneTransactions[zoneId]
	if !ok {
		e.zoneTransactions[zoneId] = []int{}
	}

	e.zoneTransactions[zoneId] = append(e.zoneTransactions[zoneId], index)

	// Planet Transaction
	_, ok = e.planetTransactions[planetId]
	if !ok {
		e.planetTransactions[planetId] = []int{}
	}

	e.planetTransactions[planetId] = append(e.planetTransactions[planetId], index)

	// TODO: Save Corporation-Planet Trade Relations level

	return nil
}

// Save contracts and existing trades between Corporation and Planets
func (e *Economy) addContract(corporationId uint64, planetId string, resource string, price float64, interval gameclock.GameTimeDuration, endTime gameclock.GameTimeDuration) error {
	e.rw.Lock()
	defer e.rw.Unlock()

	c := contract{
		corporationId: corporationId,
		planetId:      planetId,
		resource:      resource,
		price:         price,
		interval:      interval,
		endTime:       endTime,
	}

	index := len(e.contracts)
	e.contracts = append(e.contracts, c)

	_, ok := e.corporationContracts[corporationId]
	if !ok {
		e.corporationContracts[corporationId] = []int{}
	}

	e.corporationContracts[corporationId] = append(e.corporationContracts[corporationId], index)

	return nil
}

// TODO: acceptContract

// TODO: Calculate Resource Prices Projections
// func (e *Economy) calculateProjections(){
// }
