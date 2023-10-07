package corporation

import "github.com/luisya22/galactic-exchange/channel"

type CorpResponse struct {
	Err error
}

func (cg *CorpGroup) listen() {
	for i := 0; i < cg.Workers; i++ {
		go cg.worker(cg.CorpChan)
	}
}

// TODO: add WaitGroup to all the
func (cg *CorpGroup) worker(ch <-chan channel.CorpCommand) {

	for command := range ch {
		switch command.Action {
		case channel.GetSquad:
			// TODO: get corporation
			// squad, err := cg.
		}
	}
}
