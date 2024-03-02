package economy

import "github.com/luisya22/galactic-exchange/internal/gamecomm"

func (e *Economy) Listen() {
	for i := 0; i < e.Workers; i++ {
		go e.worker(e.gameChannels.CorpChannel)
	}
}

// TODO: add WaitGroup to all the workers
// TODO: Test
func (e *Economy) worker(ch <-chan gamecomm.CorpCommand) {

	for command := range ch {
		switch command.Action {
		}
	}
}
