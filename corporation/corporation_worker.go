package corporation

import "github.com/luisya22/galactic-exchange/gamecomm"

type CorpResponse struct {
	Err error
}

func (cg *CorpGroup) listen() {
	for i := 0; i < cg.Workers; i++ {
		go cg.worker(cg.CorpChan)
	}
}

// TODO: add WaitGroup to all the
func (cg *CorpGroup) worker(ch <-chan gamecomm.CorpCommand) {

	for command := range ch {
		switch command.Action {
		case gamecomm.GetSquad:
			// TODO: get corporation
			// squad, err := cg.
		}
	}
}
