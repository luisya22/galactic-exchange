package corporation

import (
	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/world"
)

func (cg *CorpGroup) listen() {
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
				break
			}

			squad, err := corp.GetSquad(command.SquadIndex)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				break
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: squad}
		case gamecomm.GetCorporation:
			corp, err := cg.findCorporationReference(command.CorporationId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				break
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: corp}
		case gamecomm.AddResourcesToBase:
			amount, err := cg.AddResources(command.CorporationId, world.Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				break
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
				break
			}

			amount, err := corp.AddResourceToSquad(command.SquadIndex, world.Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				break
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: amount}
		case gamecomm.RemoveResourcesFromSquad:
			corp, err := cg.findCorporationReference(command.CorporationId)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				break
			}

			amount, err := corp.RemoveResourcesFromSquad(command.SquadIndex, world.Resource(command.Resource), command.Amount)
			if err != nil {
				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
				break
			}

			command.ResponseChannel <- gamecomm.ChanResponse{Val: amount}
		default:
			// TODO: Handle
		}
	}
}
