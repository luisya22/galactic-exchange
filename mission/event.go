package mission

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/luisya22/galactic-exchange/channel"
)

type Event struct {
	Id        string
	MissionId string
	Time      time.Time
	Cancelled bool
	Index     int
	Execute   func(*channel.GameChannels)
}

type EventScheduler struct {
	events       map[string]*Event
	queue        EventQueue
	rw           sync.RWMutex
	gameChannels *channel.GameChannels
}

func NewEventScheduler(gameChannels *channel.GameChannels) *EventScheduler {
	return &EventScheduler{
		events:       make(map[string]*Event),
		queue:        make(EventQueue, 0),
		gameChannels: gameChannels,
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
			event.Execute(s.gameChannels)
			s.rw.Lock()
			delete(s.events, event.Id)
			s.rw.Unlock()
		} else {
			heap.Push(&s.queue, event)
		}

	}
}
