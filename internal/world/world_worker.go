package world

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

func (w *World) Listen() {
	for i := 0; i < w.Workers; i++ {
		go w.worker(w.WorldChan)
	}
}

// TODO: add WaitGroup to all the
func (w *World) worker(ch <-chan gamecomm.WorldCommand) {

	for command := range ch {
		switch command.Action {
		case gamecomm.GetPlanet:
			planet, err := w.GetPlanet(command.PlanetId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
			}

			// Return chanel
			command.ResponseChannel <- gamecomm.ChanResponse{
				Val: planet,
				Err: nil,
			}

		case gamecomm.AddResourcesToPlanet:
			amount, err := w.AddResourcesToPlanet(command.PlanetId, command.Resource, command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
			}

			command.ResponseChannel <- gamecomm.ChanResponse{
				Val: amount,
				Err: err,
			}

		case gamecomm.RemoveResourcesFromPlanet:
			amount, err := w.RemoveResourcesFromPlanet(command.PlanetId, command.Resource, command.Amount)

			command.ResponseChannel <- gamecomm.ChanResponse{
				Val: amount,
				Err: err,
			}
		case gamecomm.GetZone:
			zone, err := w.GetZone(command.ZoneId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
			}

			command.ResponseChannel <- gamecomm.ChanResponse{
				Val: zone,
				Err: nil,
			}

		default:
			command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf("error: wrong action")}

		}

		close(command.ResponseChannel)
	}
}
