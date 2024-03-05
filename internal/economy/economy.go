package economy

import (
	"sync"

	"github.com/luisya22/galactic-exchange/internal/gameclock"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/resource"
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
	corporationContracts           map[uint64][]int
	gameChannels                   gamecomm.GameChannels
	Workers                        int
	rw                             sync.RWMutex
	limit                          int
	MarketListings                 map[string][]MarketListing
	zoneMarketListingCounter       map[string]int
	zoneMutexes                    map[string]*sync.RWMutex
	resources                      map[string]resource.Resource
	gameClock                      *gameclock.GameClock
}

// type contract struct {
// 	corporationId uint64
// 	planetId      string
// 	resource      string
// 	price         float64
// 	interval      gameclock.GameTimeDuration
// 	endTime       gameclock.GameTime
// }

func NewEconomy(gameChannels gamecomm.GameChannels, resources map[string]resource.Resource, zoneIds []string, gc *gameclock.GameClock) *Economy {

	sellOffers := make(map[string][]MarketListing, len(zoneIds))
	zoneSellOfferCounter := make(map[string]int, len(zoneIds))
	zoneMutexes := make(map[string]*sync.RWMutex, len(zoneIds))

	for _, zoneId := range zoneIds {
		sellOffers[zoneId] = []MarketListing{}
		zoneSellOfferCounter[zoneId] = 0
		zoneMutexes[zoneId] = new(sync.RWMutex)
	}

	return &Economy{
		transactions:                   []transaction{},
		zoneTransactions:               make(map[string][]int),
		corporationPlanetTradeRelation: make(map[uint64]int),
		corporationContracts:           make(map[uint64][]int),
		gameChannels:                   gameChannels,
		MarketListings:                 sellOffers,
		zoneMarketListingCounter:       zoneSellOfferCounter,
		zoneMutexes:                    zoneMutexes,
		resources:                      resources,
		gameClock:                      gc,
	}
}

func (e *Economy) Run() {
	e.Listen()
}

// Save contracts and existing trades between Corporation and Planets
// func (e *Economy) addContract(corporationId uint64, planetId string, resource string, price float64, interval gameclock.GameTimeDuration, endTime gameclock.GameTime) error {
// 	e.rw.Lock()
// 	defer e.rw.Unlock()
//
// 	c := contract{
// 		corporationId: corporationId,
// 		planetId:      planetId,
// 		resource:      resource,
// 		price:         price,
// 		interval:      interval,
// 		endTime:       endTime,
// 	}
//
// 	fmt.Println(c.corporationId, c.planetId, c.resource, c.price, c.interval, c.endTime)
//
// 	index := len(e.contracts)
// 	e.contracts = append(e.contracts, c)
//
// 	_, ok := e.corporationContracts[corporationId]
// 	if !ok {
// 		e.corporationContracts[corporationId] = []int{}
// 	}
//
// 	e.corporationContracts[corporationId] = append(e.corporationContracts[corporationId], index)
//
// 	return nil
// }

// TODO: acceptContract

// TODO: Calculate Resource Prices Projections
// func (e *Economy) calculateProjections(){
// }

// TODO: Per zone marketplace
// TODO: Player and NPCs would set prices and planets would buy
// TODO: Planets would analyze and score every offer available and buy
