package mission_test

import (
	"container/heap"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/mission"
)

func createTestMissionScheduller(missions map[string]*mission.Mission, gameChannels *gamecomm.GameChannels, gc *gameclock.GameClock, mockEventScheduler *MockEventScheduler) *mission.MissionScheduler {

	return &mission.MissionScheduler{
		Missions:       missions,
		EventScheduler: mockEventScheduler,
		MissionChannel: gameChannels.MissionChannel,
		GameClock:      gc,
		GameChannels:   gameChannels,
	}

}

type MockEventScheduler struct {
	events               map[string]*mission.Event
	queue                mission.EventQueue
	gameChannels         *gamecomm.GameChannels
	missions             map[string]*mission.Mission
	gameClock            *gameclock.GameClock
	calledFunctions      map[string]int
	scheduleError        bool
	scheduleCallsToError int
	scheduleCalls        int
	updateError          bool
	rw                   sync.Mutex
}

func newMockScheduler(gameChannels *gamecomm.GameChannels, missions map[string]*mission.Mission, gc *gameclock.GameClock, scheduleError bool, scheduleCallsToError int) *MockEventScheduler {
	return &MockEventScheduler{
		events:               make(map[string]*mission.Event),
		queue:                make(mission.EventQueue, 0),
		gameChannels:         gameChannels,
		missions:             missions,
		gameClock:            gc,
		calledFunctions:      make(map[string]int),
		scheduleError:        scheduleError,
		scheduleCallsToError: scheduleCallsToError,
	}
}

func (es *MockEventScheduler) Schedule(e *mission.Event) (string, error) {

	es.scheduleCalls++

	if es.scheduleError && es.scheduleCalls == es.scheduleCallsToError {
		return "", fmt.Errorf("mock error")
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		e.Id = fmt.Sprintf("%v", len(es.events))
	}

	e.Id = uuid.String()

	es.rw.Lock()
	es.events[e.Id] = e
	heap.Push(&es.queue, e)
	es.rw.Unlock()

	es.calledFunctions["Schedule"]++

	return e.Id, nil
}

// TODO: don't love the idea of reimplementing
func (es *MockEventScheduler) UpdateEvent(eventId string, newTime gameclock.GameTime, cancelled bool) error {

	es.calledFunctions["UpdateEvent"]++

	var event *mission.Event
	var ok bool

	es.rw.Lock()
	if event, ok = es.events[eventId]; !ok {
		return fmt.Errorf("error: event not found %v", eventId)
	}

	es.queue.Update(event, newTime, cancelled)
	es.rw.Unlock()

	return nil
}

func (es *MockEventScheduler) Run() {

}

type testError struct {
	shouldError bool
	errorStr    string
}

type worldErrors struct {
	getPlanetError        testError
	addResourcesPlanet    testError
	removeResourcesPlanet testError
}

func listenWordlWorker(t *testing.T, errors worldErrors, gameChannels *gamecomm.GameChannels, wg *sync.WaitGroup) {
	t.Helper()

	go func(ch chan gamecomm.WorldCommand) {
		for command := range ch {
			switch command.Action {
			case gamecomm.GetPlanet:

				if errors.getPlanetError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.getPlanetError.errorStr)}
					continue
				}

				planet := gamecomm.Planet{}

				command.ResponseChannel <- gamecomm.ChanResponse{
					Val: planet,
					Err: nil,
				}

			case gamecomm.AddResourcesToPlanet:

				if errors.addResourcesPlanet.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.addResourcesPlanet.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{
					Val: 0,
					Err: nil,
				}

			case gamecomm.RemoveResourcesFromPlanet:

				if errors.removeResourcesPlanet.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.removeResourcesPlanet.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{
					Val: 0,
					Err: nil,
				}
			default:
				command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf("error: wrong action")}

			}
		}

		wg.Done()
	}(gameChannels.WorldChannel)
}

type corporationErrors struct {
	getSquadError                 testError
	getCorporationError           testError
	addResourcesToBaseError       testError
	removeResourcesFromBaseError  testError
	addResourcesToSquad           testError
	removeResourcesFromSquadError testError
	addCreditsError               testError
	removeCreditsError            testError
}

func listenCorporationWorker(t *testing.T, errors corporationErrors, gamechannels *gamecomm.GameChannels, wg *sync.WaitGroup) {
	t.Helper()

	go func(ch chan gamecomm.CorpCommand) {

		for command := range ch {
			switch command.Action {
			case gamecomm.GetSquad:

				if errors.getSquadError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.getSquadError.errorStr)}
					continue
				}

				squad := gamecomm.Squad{}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: squad}
			case gamecomm.GetCorporation:

				if errors.getCorporationError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.getCorporationError.errorStr)}
					continue
				}

				corp := gamecomm.Corporation{}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: corp}
			case gamecomm.AddResourcesToBase:

				if errors.addResourcesToBaseError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.addResourcesToBaseError.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: 1}
			case gamecomm.RemoveResourcesFromBase:

				if errors.removeResourcesFromBaseError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.removeResourcesFromBaseError.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: 1}
			case gamecomm.AddResourcesToSquad:

				if errors.addResourcesToSquad.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.addResourcesToSquad.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: 1}
			case gamecomm.RemoveResourcesFromSquad:

				if errors.removeResourcesFromSquadError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.removeResourcesFromSquadError.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: 1}
			case gamecomm.AddCredits:

				if errors.addCreditsError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.addCreditsError.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: 1}
			case gamecomm.RemoveCredits:

				if errors.removeCreditsError.shouldError {
					command.ResponseChannel <- gamecomm.ChanResponse{Err: fmt.Errorf(errors.removeCreditsError.errorStr)}
					continue
				}

				command.ResponseChannel <- gamecomm.ChanResponse{Val: 1}

			default:
				// TODO: Handle
			}
		}

		wg.Done()
	}(gamechannels.CorpChannel)
}

// type EventScheduler interface {
// 	Schedule(*Event)
// 	UpdateEvent(string, gameclock.GameTime, bool) error
// 	Run()
// }

// type DefaultEventScheduler struct {
// 	events       map[string]*Event
// 	queue        EventQueue
// 	rw           sync.RWMutex
// 	gameChannels *gamecomm.GameChannels
// 	missions     map[string]*Mission
// 	gameClock    *gameclock.GameClock
// }

// func NewMissionScheduler(gameChannels *gamecomm.GameChannels, gc *gameclock.GameClock) *MissionScheduler {
//
// 	missions := make(map[string]*Mission, 0)
// 	eventScheduler := NewEventScheduler(gameChannels, missions, gc)
//
// 	return &MissionScheduler{
// 		missions:       missions,
// 		eventScheduler: eventScheduler,
// 		missionChannel: gameChannels.MissionChannel,
// 		gameClock:      gc,
// 		gameChannels:   gameChannels,
// 	}
// }

// type MissionScheduler struct {
// 	missions       map[string]*Mission
// 	eventScheduler *EventScheduler
// 	missionChannel chan gamecomm.MissionCommand
// 	RW             sync.RWMutex
// 	gameClock      *gameclock.GameClock
// }
