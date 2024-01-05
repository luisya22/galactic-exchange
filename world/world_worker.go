package world

import (
	"github.com/luisya22/galactic-exchange/gamecomm"
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
				Val: planet.copy(),
				Err: nil,
			}

			break
		case gamecomm.AddResourcesToPlanet:
			amount, err := w.AddResourcesToPlanet(command.PlanetId, Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
			}

			command.ResponseChannel <- gamecomm.ChanResponse{
				Val: amount,
				Err: err,
			}

			break
		case gamecomm.RemoveResourcesFromPlanet:
			amount, err := w.RemoveResourcesFromPlanet(command.PlanetId, Resource(command.Resource), command.Amount)

			command.ResponseChannel <- gamecomm.ChanResponse{
				Val: amount,
				Err: err,
			}

			break

		}
	}
}
