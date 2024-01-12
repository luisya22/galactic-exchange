package corporation

import (
	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/world"
)

func (cg *CorpGroup) Listen() {
	for i := 0; i < cg.Workers; i++ {
		go cg.worker(cg.CorpChan)
	}
}

// TODO: add WaitGroup to all the workers
// TODO: Test
func (cg *CorpGroup) worker(ch <-chan gamecomm.CorpCommand) {

	for command := range ch {
		switch command.Action {
		case gamecomm.GetSquad:
			corp, err := cg.findCorporationReference(command.CorporationId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			squad, err := corp.GetSquad(command.SquadIndex)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: squad}
		case gamecomm.GetCorporation:
			corp, err := cg.findCorporation(command.CorporationId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: corp}
		case gamecomm.AddResourcesToBase:
			amount, err := cg.AddResources(command.CorporationId, world.Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: amount}
		case gamecomm.RemoveResourcesFromBase:
			amount, err := cg.RemoveResources(command.CorporationId, world.Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: amount}
		case gamecomm.AddResourcesToSquad:
			corp, err := cg.findCorporationReference(command.CorporationId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			amount, err := corp.AddResourceToSquad(command.SquadIndex, world.Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: amount}
		case gamecomm.RemoveResourcesFromSquad:
			corp, err := cg.findCorporationReference(command.CorporationId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			amount, err := corp.RemoveResourcesFromSquad(command.SquadIndex, world.Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: amount}
		case gamecomm.RemoveAllResourcesFromSquad:
			corp, err := cg.findCorporationReference(command.CorporationId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			amount, err := corp.RemoveAllResourcesFromSquad(command.SquadIndex, world.Resource(command.Resource))
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}
		case gamecomm.AddCredits:
			credits, err := cg.AddCredits(command.CorporationId, command.AmountDecimal)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: credits}
		case gamecomm.RemoveCredits:
			credits, err := cg.RemoveCredits(command.CorporationId, command.AmountDecimal)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				continue
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: credits}

		default:
			// TODO: Handle
		}
	}
}
