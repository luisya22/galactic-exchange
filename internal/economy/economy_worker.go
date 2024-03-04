package economy

import (
	"log"

	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

func (e *Economy) Listen() {
	for i := 0; i < e.Workers; i++ {
		go e.worker(e.gameChannels.EconomyChannel)
	}
}

// TODO: add WaitGroup to all the workers
// TODO: Test
func (e *Economy) worker(ch <-chan gamecomm.EconomyCommand) {

	for command := range ch {
		switch command.Action {
		case gamecomm.AddMarketListing:
			so := MarketListing{
				ResourceName:  command.Resource,
				Amount:        command.Amount,
				Price:         command.Price,
				CorporationId: command.CorporationId,
			}

			id, err := e.addMarketListing(command.ZoneId, so)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				return
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: id}

		case gamecomm.BuyMarketListing:
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

			command.ResponseChannel <- gamecomm.ChanResponse{Val: marketListings}
		case gamecomm.EditMarketListingPrice:
			err := e.editPrice(command.ZoneId, command.MarketListingId, command.CorporationId, command.Price)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				close(command.ResponseChannel)
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: "OK"}
		}

		close(command.ResponseChannel)
	}
}
