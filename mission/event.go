package mission

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/luisya22/galactic-exchange/gamecomm"
)

type Event struct {
	Id        string
	MissionId string
	Time      time.Time
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
}

func NewEventScheduler(gameChannels *gamecomm.GameChannels, missions map[string]*Mission) *EventScheduler {
	return &EventScheduler{
		events:       make(map[string]*Event),
		queue:        make(EventQueue, 0),
		gameChannels: gameChannels,
		missions:     missions,
	}
}

func (s *EventScheduler) Schedule(e *Event) {
	s.rw.Lock()
	s.events[e.Id] = e
	heap.Push(&s.queue, e)
	s.rw.Unlock()
}

func (s *EventScheduler) UpdateEvent(eventId string, newTime time.Time, cancelled bool) error {
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

		now := time.Now()
		if now.After(event.Time) {
			mission := s.missions[event.MissionId]
			event.Execute(mission, s.gameChannels)
			s.rw.Lock()
			delete(s.events, event.Id)
			s.rw.Unlock()
		} else {
			heap.Push(&s.queue, event)
		}

	}
}
