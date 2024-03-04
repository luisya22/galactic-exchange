package economy

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

// MarketRequest

type MarketListing struct {
	Id              string
	ResourceName    string
	Amount          int
	RemainingAmount int
	Price           float64
	CorporationId   uint64
}

func (e *Economy) addMarketListing(zoneId string, so MarketListing) (string, error) {
	if _, ok := e.resources[so.ResourceName]; !ok {
		return "", fmt.Errorf("error: resource doesn't exist")
	}

	if so.Amount <= 0 {
		return "", fmt.Errorf("error: amount should be greater than zero")
	}

	if so.Price <= 0 {
		return "", fmt.Errorf("error: price should be greater than zero")
	}

	resChan := make(chan gamecomm.ChanResponse)

	e.gameChannels.CorpChannel <- gamecomm.CorpCommand{
		Action:          gamecomm.GetCorporation,
		CorporationId:   so.CorporationId,
		ResponseChannel: resChan,
	}

	res := <-resChan
	if res.Err != nil {
		return "", fmt.Errorf("error: incorrect corporation id")
	}

	e.rw.RLock()
	mutex, ok := e.zoneMutexes[zoneId]
	e.rw.RUnlock()
	if !ok {
		return "", fmt.Errorf("error: no mutex found for zone ID '%s'", zoneId)
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Generate Id
	e.zoneMarketListingCounter[zoneId]++
	num := e.zoneMarketListingCounter[zoneId]
	id := fmt.Sprintf("%v-%d", zoneId, num)

	so.Id = id
	e.MarketListings[zoneId] = append(e.MarketListings[zoneId], so)

	return so.Id, nil
}

func (e *Economy) getZoneMarketListings(zoneId string) ([]MarketListing, error) {
	e.rw.RLock()
	mutex, ok := e.zoneMutexes[zoneId]
	e.rw.RUnlock()
	if !ok {
		return []MarketListing{}, fmt.Errorf("error: no mutex found for zone ID '%s'", zoneId)
	}

	mutex.RLock()
	defer mutex.RUnlock()

	marketListings, ok := e.MarketListings[zoneId]
	if !ok {
		return []MarketListing{}, nil
	}

	return marketListings, nil
}

func (e *Economy) getMarketListing(zoneId string, listingId string) (MarketListing, error) {
	e.rw.RLock()
	mutex, ok := e.zoneMutexes[zoneId]
	e.rw.RUnlock()
	if !ok {
		return MarketListing{}, fmt.Errorf("error: no mutex found for zone ID '%s'", zoneId)
	}

	mutex.RLock()
	defer mutex.RUnlock()

	marketListings, ok := e.MarketListings[zoneId]
	if !ok {
		return MarketListing{}, fmt.Errorf("error: no market listing found for zone ID '%s'", zoneId)
	}

	for _, ml := range marketListings {
		if ml.Id == listingId {
			return ml, nil
		}
	}

	return MarketListing{}, fmt.Errorf("error: listing with ID '%s' not found", listingId)
}

func (e *Economy) removeAmount(zoneId string, listingId string, amount int) (int, error) {
	e.rw.RLock()
	mutex, ok := e.zoneMutexes[zoneId]
	e.rw.Unlock()
	if !ok {
		return 0, fmt.Errorf("error: no mutex found for zone ID '%s'", zoneId)
	}

	mutex.Lock()
	defer mutex.Unlock()

	marketListings, ok := e.MarketListings[zoneId]
	if !ok {
		return 0, fmt.Errorf("error: no market listing found for zone ID '%s'", zoneId)
	}

	selectedIndex := -1
	for i, ml := range marketListings {
		if ml.Id == listingId {
			selectedIndex = i
		}
	}

	if selectedIndex == -1 {
		return 0, fmt.Errorf("error: listing not found with ID '%s'", listingId)
	}

	remainingAmount := e.MarketListings[zoneId][selectedIndex].Amount
	if remainingAmount < amount {
		amount = remainingAmount
	}

	e.MarketListings[zoneId][selectedIndex].Amount -= amount

	if e.MarketListings[zoneId][selectedIndex].Amount == 0 {
		e.MarketListings[zoneId] = append(marketListings[:selectedIndex], marketListings[selectedIndex+1:]...)
	}

	return amount, nil
}

// func (e *Economy) removeMarketListing(zoneId string, listingId string) error {
//
// 	e.rw.RLock()
// 	mutex, ok := e.zoneMutexes[zoneId]
// 	e.rw.Unlock()
// 	if !ok {
// 		return fmt.Errorf("error: no mutex found for zone ID '%s'", zoneId)
// 	}
//
// 	mutex.Lock()
// 	defer mutex.Unlock()
//
// 	marketListings, ok := e.MarketListings[zoneId]
// 	if !ok {
// 		return fmt.Errorf("error: no market listing found for zone ID '%s'", zoneId)
// 	}
//
// 	selectedIndex := -1
// 	for i, ml := range marketListings {
// 		if ml.Id == listingId {
// 			selectedIndex = i
// 		}
// 	}
//
// 	e.MarketListings[zoneId] = append(marketListings[:selectedIndex], marketListings[selectedIndex+1:]...)
//
// 	return nil
// }

func (e *Economy) editPrice(zoneId string, listingId string, corporationId uint64, price float64) error {
	e.rw.RLock()
	mutex, ok := e.zoneMutexes[zoneId]
	e.rw.Unlock()
	if !ok {
		return fmt.Errorf("error: no mutex found for zone ID '%s'", zoneId)
	}

	mutex.Lock()
	defer mutex.Unlock()

	marketListings, ok := e.MarketListings[zoneId]
	if !ok {
		return fmt.Errorf("error: no market listing found for zone ID '%s'", zoneId)
	}

	selectedIndex := -1
	for i, ml := range marketListings {
		if ml.Id == listingId {
			selectedIndex = i
		}
	}

	listing := e.MarketListings[zoneId][selectedIndex]

	if listing.CorporationId != corporationId {
		return fmt.Errorf("error: you can not edit listing with ID '%s'", listingId)
	}

	e.MarketListings[zoneId][selectedIndex].Price = price

	return nil
}
