package mission

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
)

type Event struct {
	Id        string
	MissionId string
	Time      gameclock.GameTime
	Cancelled bool
	Index     int
	Execute   func(*Mission, *gamecomm.GameChannels)
}

type EventScheduler struct {
	events       map[string]*Event
	queue        EventQueue
	rw           sync.RWMutex
	gameChannels *gamecomm.GameChannels
	missions     map[string]*Mission
	gameClock    *gameclock.GameClock
}

func NewEventScheduler(gameChannels *gamecomm.GameChannels, missions map[string]*Mission, gc *gameclock.GameClock) *EventScheduler {
	return &EventScheduler{
		events:       make(map[string]*Event),
		queue:        make(EventQueue, 0),
		gameChannels: gameChannels,
		missions:     missions,
		gameClock:    gc,
	}
}

func (s *EventScheduler) Schedule(e *Event) {
	s.rw.Lock()
	s.events[e.Id] = e
	heap.Push(&s.queue, e)
	s.rw.Unlock()
}

func (s *EventScheduler) UpdateEvent(eventId string, newTime gameclock.GameTime, cancelled bool) error {
	var event *Event
	var ok bool

	s.rw.Lock()
	if event, ok = s.events[eventId]; !ok {
		return fmt.Errorf("error: event not found %v", eventId)
	}

	s.queue.Update(event, newTime, cancelled)
	s.rw.Unlock()

	return nil
}

func (s *EventScheduler) Run() {
	for {
		if len(s.queue) == 0 {
			continue
		}

		event := heap.Pop(&s.queue).(*Event)

		if event.Cancelled {
			s.rw.Lock()
			delete(s.events, event.Id)
			s.rw.Unlock()
			continue
		}

		now := s.gameClock.GetCurrentTime()
		if now.After(event.Time) {
			mission := s.missions[event.MissionId]
			go event.Execute(mission, s.gameChannels)
			s.rw.Lock()
			delete(s.events, event.Id)
			s.rw.Unlock()
		} else {
			heap.Push(&s.queue, event)
		}

	}
}
