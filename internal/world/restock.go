package world

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/internal/economy"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

type listingScore struct {
	score float64
	index int
}

// TODO: After doing buy move everything to it's own package

func (w *World) basicSupplyRestock(planetId string, resourceName string, monthlyProduction int, actualStock int) {
	if actualStock/30 < monthlyProduction {
		// TODO: Buy
		w.RW.RLock()
		defer w.RW.RUnlock()

		planet := w.Planets[planetId]

		wantsToBuy := monthlyProduction * w.randomInt(1, 3)
		planet.buyResource(resourceName, wantsToBuy, w.economyChan)

		fmt.Printf("%v wants to buy %v %v\n", planetId, wantsToBuy, resourceName)

	}
}

func (w *World) processResourceConsumption(planetId string, resourceName string, minConsumption int, maxConsumption int) {
	quantity := w.randomInt(minConsumption, maxConsumption)
	remaning, _ := w.DepletePlanetResource(planetId, resourceName, quantity)

	weeklyConsumption := quantity * 7

	w.restockResources(planetId, resourceName, weeklyConsumption, remaning)
}

func classifyResourceLevel(weeklyConsumption int, totalStorage int) string {
	if totalStorage == 0 {
		return "depleted"
	}

	consumptionPercentage := (weeklyConsumption / totalStorage) * 100

	switch {
	case consumptionPercentage < 10:
		return "high"
	case consumptionPercentage >= 10 && consumptionPercentage < 50:
		return "medium"
	case consumptionPercentage >= 50 && consumptionPercentage < 75:
		return "low"
	default:
		return "depleted"
	}
}

func purchaseProbability(resourceLevel string) float64 {
	switch resourceLevel {
	case "high":
		return 0.10
	case "medium":
		return 0.50
	case "low":
		return 0.75
	case "depleted":
		return 1.00
	default:
		return 0.00
	}
}

func (w *World) restockResources(planetId string, resource string, weeklyConsumption int, totalStorage int) {
	resourceLevel := classifyResourceLevel(weeklyConsumption, totalStorage)
	purchaseProb := purchaseProbability(resourceLevel)

	randomFloat := w.RandomNumber.Float64()

	wantsToBuy := w.randomInt(weeklyConsumption, weeklyConsumption*4)

	if randomFloat <= purchaseProb {
		// TODO: Make purchase
		w.RW.RLock()
		defer w.RW.RUnlock()

		planet := w.Planets[planetId]
		planet.buyResource(resource, wantsToBuy, w.economyChan)

	}
}

func (planet *Planet) buyResource(resource string, amount int, economyChan chan gamecomm.EconomyCommand) {
	planet.RW.RLock()
	defer planet.RW.RUnlock()

	// Get Market Price
	resChan := make(chan gamecomm.ChanResponse)
	command := gamecomm.EconomyCommand{
		Action:          gamecomm.GetMarketPrice,
		ZoneId:          planet.ZoneId,
		ResponseChannel: resChan,
	}

	economyChan <- command

	res := <-resChan
	if res.Err != nil {
		return
	}

	marketValue := res.Val.(float64)

	// Get all offers
	resChan = make(chan gamecomm.ChanResponse)
	command = gamecomm.EconomyCommand{
		Action:          gamecomm.GetMarketListingsByResource,
		Resource:        resource,
		ResponseChannel: resChan,
	}

	economyChan <- command
	res = <-resChan
	if res.Err != nil {
		return
	}

	marketListings := res.Val.([]economy.MarketListing)

	// TODO: LUIS HERE EVALUATE

	// Evaluate and score each one

	// TODO: attach in some way listingScore to MarketListing.Id
	listingScores := []listingScore{}
	for i, ml := range marketListings {
		score := scoreListing(ml, amount, marketValue)

		ls := listingScore{
			score: score,
			index: i,
		}

		listingScores = insertSorted(listingScores, ls)
	}

	for i := range listingScores {

		ml := marketListings[i]

		resChan := make(chan gamecomm.ChanResponse)
		command := gamecomm.EconomyCommand{
			Action:          gamecomm.BuyMarketListing,
			ZoneId:          planet.ZoneId,
			MarketListingId: ml.Id,
			ResponseChannel: resChan,
		}

		economyChan <- command
		res = <-resChan
		if res.Err != nil {
			return
		}

		amount -= res.Val.(int)
		if amount <= 0 {
			break
		}
	}

}

func scoreListing(marketListing economy.MarketListing, amountToBuy int, marketValue float64) float64 {
	priceScore := marketValue / marketListing.Price

	fullAmountScore := 0.0
	if marketListing.Amount >= amountToBuy {
		fullAmountScore = 1
	}

	priceWeight := 0.7
	amountWeight := 0.3

	score := (priceWeight * priceScore) + (amountWeight * fullAmountScore)
	return score
}

func binarySearch(listingScores []listingScore, ls listingScore) int {
	low, high := 0, len(listingScores)
	for low < high {
		mid := low + (high-low)/2
		if listingScores[mid].score < ls.score {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low
}

func insertSorted(listingScores []listingScore, ls listingScore) []listingScore {
	index := binarySearch(listingScores, ls)

	listingScores = append(listingScores, listingScore{})

	copy(listingScores[index+1:], listingScores[index:])

	listingScores[index] = ls

	return listingScores
}

// TODO: planets should analyze their resource scarcity
