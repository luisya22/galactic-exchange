package world

import "github.com/luisya22/galactic-exchange/gamecomm"

type WorldResponse struct {
	Planet Planet
	Amount int
	Err    error
}

func (w *World) listen() {
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
				command.ResponseChannel <- WorldResponse{Err: err}
			}

			// Return chanel
			command.ResponseChannel <- WorldResponse{
				Planet: planet.copy(),
				Err:    nil,
			}
		case gamecomm.AddResourcesToPlanet:
			amount, err := w.AddResourcesToPlanet(command.PlanetId, Resource(command.Resource), command.Amount)

			command.ResponseChannel <- WorldResponse{
				Amount: amount,
				Err:    err,
			}
		case gamecomm.RemoveResourcesFromPlanet:
			amount, err := w.RemoveResourcesFromPlanet(command.PlanetId, Resource(command.Resource), command.Amount)

			command.ResponseChannel <- WorldResponse{
				Amount: amount,
				Err:    err,
			}

		}
	}
}
