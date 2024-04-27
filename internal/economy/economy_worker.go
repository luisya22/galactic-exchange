package economy

import (
	"fmt"
	"log"
	"time"

	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/resource"
)

func (e *Economy) listen() {
	for i := 0; i < e.Workers; i++ {
		go e.worker(e.gameChannels.EconomyChannel)
	}
}

func (e *Economy) addRandomMarketListings(resources []resource.Resource, zoneIds []string, economyChannel chan gamecomm.EconomyCommand) {

	for x := 0; x < 100; x++ {
		for _, z := range zoneIds {
			for _, r := range resources {
				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.EconomyCommand{
					Action:          gamecomm.AddMarketListing,
					Resource:        r.Name,
					Amount:          100_000,
					Price:           100,
					CorporationId:   1,
					ZoneId:          z,
					ResponseChannel: resChan,
				}

				economyChannel <- command

				res := <-resChan
				if res.Err != nil {
					fmt.Println(res.Err.Error())
				}
			}
		}

		time.Sleep(5 * time.Second)
	}

}

// TODO: add WaitGroup to all the workers
// TODO: Test
func (e *Economy) worker(ch <-chan gamecomm.EconomyCommand) {
	for command := range ch {

		listingTime := e.gameClock.GetCurrentTime()

		switch command.Action {
		case gamecomm.AddMarketListing:
			so := MarketListing{
				ResourceName:  command.Resource,
				Amount:        command.Amount,
				Price:         command.Price,
				CorporationId: command.CorporationId,
				ListTime:      listingTime,
			}

			id, err := e.addMarketListing(command.ZoneId, so)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				return
			}

			// if _, ok := e.zoneAnalytics[command.ZoneId]; ok {
			// 	e.zoneAnalytics[command.ZoneId].updateListingAmount(command.Resource, listingTime)
			// 	e.zoneAnalytics[command.ZoneId].updateListingVolume(command.Resource, command.Amount, listingTime)
			// }

			command.ResponseChannel <- gamecomm.ChanResponse{Val: id}

		case gamecomm.BuyMarketListing:
			fmt.Println(command)
			marketListing, err := e.getMarketListing(command.ZoneId, command.MarketListingId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				continue
			}

			amount, err := e.removeAmount(command.ZoneId, command.MarketListingId, command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				continue
			}

			// if _, ok := e.zoneAnalytics[command.ZoneId]; ok {
			// 	e.zoneAnalytics[command.ZoneId].updateSalesAmount(command.Resource, listingTime)
			// 	e.zoneAnalytics[command.ZoneId].updateSalesVolume(command.Resource, command.Amount, listingTime)
			// }

			command.ResponseChannel <- gamecomm.ChanResponse{Val: amount}

			err = e.addTransaction(
				command.ZoneId,
				command.BuyerPlanetId,
				command.CorporationId,
				command.Resource,
				marketListing.Price,
				e.gameClock.GetCurrentTime(),
			)
			if err != nil {
				log.Println(err)
			}
		case gamecomm.GetMarketListings:
			marketListings, err := e.getZoneMarketListings(command.ZoneId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: *marketListings}
		case gamecomm.EditMarketListingPrice:
			err := e.editPrice(command.ZoneId, command.MarketListingId, command.CorporationId, command.Price)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: "OK"}
		case gamecomm.GetMarketListingsByResource:
			marketListings, err := e.getZoneMarketListingsByResource(command.ZoneId, command.Resource)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: *marketListings}
		case gamecomm.GetMarketPrice:
			marketPrice, err := e.getZoneResourceMarketPrice(command.ZoneId, command.Resource)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: marketPrice}
		}

		close(command.ResponseChannel)
	}
}
