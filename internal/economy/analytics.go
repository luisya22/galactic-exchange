package economy

import (
	"sync"

	"github.com/luisya22/galactic-exchange/internal/gameclock"
	"github.com/luisya22/galactic-exchange/internal/resource"
)

type zoneAnalytics map[string]*analytics
type resourceVolume map[string]int
type transactionFrequency map[string]int
type avgListingDuration map[string]float64
type dailyPrices map[gameclock.GameTime]float64

type analytics struct {
	salesVolume        map[gameclock.GameTime]resourceVolume
	salesAmount        map[gameclock.GameTime]transactionFrequency
	avgListingDuration map[gameclock.GameTime]avgListingDuration
	listingVolume      map[gameclock.GameTime]resourceVolume
	listingAmount      map[gameclock.GameTime]transactionFrequency
	rw                 sync.RWMutex
	historicPrices     map[string]dailyPrices
}

func newAnalytics() *analytics {
	return &analytics{
		salesVolume:        make(map[gameclock.GameTime]resourceVolume),
		salesAmount:        make(map[gameclock.GameTime]transactionFrequency),
		avgListingDuration: make(map[gameclock.GameTime]avgListingDuration),
		listingVolume:      make(map[gameclock.GameTime]resourceVolume),
		listingAmount:      make(map[gameclock.GameTime]transactionFrequency),
		historicPrices:     make(map[string]dailyPrices),
	}
}

func (a *analytics) updateListingVolume(resource string, quantity int, listingTime gameclock.GameTime) {
	a.rw.Lock()
	defer a.rw.Unlock()

	if _, ok := a.listingVolume[listingTime.StartOfDay()]; !ok {
		a.listingVolume[listingTime.StartOfDay()] = make(resourceVolume)
	}

	if _, ok := a.listingVolume[listingTime.StartOfDay()][resource]; !ok {
		a.listingVolume[listingTime.StartOfDay()][resource] = 0
	}

	a.listingVolume[listingTime.StartOfDay()][resource] += quantity
}

func (a *analytics) updateListingAmount(resource string, listingTime gameclock.GameTime) {
	a.rw.Lock()
	defer a.rw.Unlock()

	if _, ok := a.listingAmount[listingTime.StartOfDay()]; !ok {
		a.listingAmount[listingTime.StartOfDay()] = make(transactionFrequency)
	}

	if _, ok := a.listingAmount[listingTime.StartOfDay()][resource]; !ok {
		a.listingAmount[listingTime.StartOfDay()][resource] = 0
	}

	a.listingAmount[listingTime.StartOfDay()][resource] += 1
}

func (a *analytics) updateSalesVolume(resource string, quantity int, saleTime gameclock.GameTime) {
	a.rw.Lock()
	defer a.rw.Unlock()

	if _, ok := a.salesVolume[saleTime.StartOfDay()]; !ok {
		a.salesVolume[saleTime.StartOfDay()] = make(resourceVolume)
	}

	if _, ok := a.salesVolume[saleTime.StartOfDay()][resource]; !ok {
		a.salesVolume[saleTime.StartOfDay()][resource] = 0
	}

	a.salesVolume[saleTime.StartOfDay()][resource] += quantity
}

func (a *analytics) updateSalesAmount(resource string, saleTime gameclock.GameTime) {
	a.rw.Lock()
	defer a.rw.Unlock()

	if _, ok := a.salesAmount[saleTime.StartOfDay()]; !ok {
		a.salesAmount[saleTime.StartOfDay()] = make(transactionFrequency)
	}

	if _, ok := a.salesAmount[saleTime.StartOfDay()][resource]; !ok {
		a.salesAmount[saleTime.StartOfDay()][resource] = 0
	}

	a.salesAmount[saleTime.StartOfDay()][resource]++
}

func (a *analytics) calculateDailySupply(resource string, day gameclock.GameTime) int {
	// a.rw.Lock()
	// defer a.rw.Unlock()

	dayStart := day.StartOfDay()
	totalSupply := 0

	if volumes, ok := a.listingVolume[dayStart]; ok {
		totalSupply = volumes[resource]
	}

	return totalSupply
}

func (a *analytics) calculateDailyDemand(resource string, day gameclock.GameTime) int {
	// a.rw.Lock()
	// defer a.rw.Unlock()

	dayStart := day.StartOfDay()
	totalDemand := 0

	if volumes, ok := a.salesVolume[dayStart]; ok {
		totalDemand = volumes[resource]
	}

	return totalDemand
}

func (a *analytics) updateItemPrices(resources map[string]resource.Resource, prices resourcePrices, day gameclock.GameTime) resourcePrices {
	// a.rw.Lock()
	// defer a.rw.Unlock()
	startOfDay := day.StartOfDay()
	for name := range resources {
		dailySupply := a.calculateDailySupply(name, startOfDay)
		dailyDemand := a.calculateDailyDemand(name, startOfDay)

		var priceAdjustment float64
		if dailyDemand > dailySupply {
			priceAdjustment = 1.01
		} else if dailyDemand < dailySupply {
			priceAdjustment = 0.99
		} else {
			priceAdjustment = 1.00
		}

		currentPrice := prices[name]
		updatedPrice := currentPrice * priceAdjustment
		prices[name] = updatedPrice
	}

	return prices
}

func (a *analytics) storeHistoricPrices(resources map[string]resource.Resource, prices resourcePrices, day gameclock.GameTime) {
	// a.rw.Lock()
	// defer a.rw.Unlock()

	startOfDay := day.StartOfDay()

	for name := range resources {
		price := prices[name]

		if _, ok := a.historicPrices[name]; !ok {
			a.historicPrices[name] = make(dailyPrices)
		}

		if _, ok := a.historicPrices[name][startOfDay]; !ok {
			a.historicPrices[name][startOfDay] = 0
		}

		a.historicPrices[name][startOfDay] = price
	}
}
