package mission

import (
	"container/heap"
	"fmt"
	"sync"
	"testing"

	"github.com/luisya22/galactic-exchange/internal/assert"
	"github.com/luisya22/galactic-exchange/internal/gameclock"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

func TestSchedule(t *testing.T) {
	// Validate that the event is added to map and heap
	// Validate that the event is added on the correct order

	tests := []struct {
		name   string
		events []*Event
		wants  []*Event
	}{
		{
			name: "It Schedules Event",
			events: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      5,
					Cancelled: false,
					Index:     0,
				},
			},
			wants: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      5,
					Cancelled: false,
					Index:     0,
				},
			},
		},
		{
			name: "It Schedules Ordered Events In Order Based On Time",
			events: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      1,
					Cancelled: false,
					Index:     0,
				},
				{
					Id:        "Event-2",
					MissionId: "Mission-1",
					Time:      2,
					Cancelled: false,
					Index:     0,
				},
			},
			wants: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      1,
					Cancelled: false,
					Index:     0,
				},
				{
					Id:        "Event-2",
					MissionId: "Mission-1",
					Time:      2,
					Cancelled: false,
					Index:     0,
				},
			},
		},
		{
			name: "It Schedules Unordered Events In Order Based On Time",
			events: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      17,
					Cancelled: false,
					Index:     0,
				},
				{
					Id:        "Event-2",
					MissionId: "Mission-1",
					Time:      2,
					Cancelled: false,
					Index:     0,
				},
			},
			wants: []*Event{
				{
					Id:        "Event-2",
					MissionId: "Mission-1",
					Time:      2,
					Cancelled: false,
					Index:     0,
				},
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      17,
					Cancelled: false,
					Index:     0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			eventScheduler := &DefaultEventScheduler{
				events:    make(map[string]*Event),
				queue:     make(EventQueue, 0),
				missions:  make(map[string]*Mission),
				gameClock: gameclock.NewGameClock(0, 1),
			}

			eventScheduler.idGenerator = func(e *Event) error {
				return nil
			}

			for _, e := range tt.events {
				_, err := eventScheduler.Schedule(e)
				if err != nil {
					t.Fatal(err)
				}
			}

			// validate events are on map
			for _, e := range tt.events {
				assert.MapContains(t, eventScheduler.events, e.Id, e)
			}

			// validate events are on heap and in order
			for index, e := range tt.wants {

				queueEvent := eventScheduler.queue[index]
				assert.Equal(t, queueEvent.Id, e.Id)
			}
		})
	}

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numEvents := 10

		eventScheduler := &DefaultEventScheduler{
			events:    make(map[string]*Event),
			queue:     make(EventQueue, 0),
			missions:  make(map[string]*Mission),
			gameClock: gameclock.NewGameClock(0, 1),
		}

		eventScheduler.idGenerator = func(e *Event) error {
			return nil
		}

		var events []*Event

		for i := numEvents - 1; i >= 0; i-- {
			event := &Event{
				Id:        fmt.Sprintf("Event-%v", i),
				MissionId: "Mission-1",
				Time:      gameclock.GameTime((i + 1) * 2),
				Cancelled: false,
			}

			events = append(events, event)
		}

		for _, event := range events {
			wg.Add(1)
			go func(e *Event, eventScheduler *DefaultEventScheduler) {
				defer wg.Done()

				_, err := eventScheduler.Schedule(e)
				if err != nil {
					t.Errorf(err.Error())
				}

			}(event, eventScheduler)
		}

		wg.Wait()

		// validate events are on map
		for _, e := range events {
			assert.MapContains(t, eventScheduler.events, e.Id, e)
		}

		// validate events are on heap and in order
		for i := numEvents - 1; i >= 0; i-- {
			queueEvent := heap.Pop(&eventScheduler.queue).(*Event)
			assert.Equal(t, queueEvent, events[i])
		}
	})
}

func TestRun(t *testing.T) {

	type eventMsg struct {
		msg  string
		time gameclock.GameTime
	}

	msgChan := make(chan eventMsg)

	tests := []struct {
		name   string
		events []*Event
		wants  []eventMsg
	}{
		{
			name: "Run Event",
			events: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      1,
					Cancelled: false,
					Execute: func(m *Mission, gc *gamecomm.GameChannels) {
						msgChan <- eventMsg{
							msg: "Hello there",
						}
					},
				},
			},
			wants: []eventMsg{
				{
					msg:  "Hello there",
					time: 1,
				},
			},
		},
		{
			name: "Run Two Events Scheduled In Order",
			events: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      1,
					Cancelled: false,
					Execute: func(m *Mission, gc *gamecomm.GameChannels) {
						msgChan <- eventMsg{
							msg: "Message 1",
						}
					},
				},
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      2,
					Cancelled: false,
					Execute: func(m *Mission, gc *gamecomm.GameChannels) {
						msgChan <- eventMsg{
							msg: "Message 2",
						}
					},
				},
			},
			wants: []eventMsg{
				{
					msg:  "Message 1",
					time: 1,
				},
				{
					msg:  "Message 2",
					time: 2,
				},
			},
		},
		{
			name: "Run Two Events Scheduled Out Of Order",
			events: []*Event{
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      3,
					Cancelled: false,
					Execute: func(m *Mission, gc *gamecomm.GameChannels) {
						msgChan <- eventMsg{
							msg: "Message 1",
						}
					},
				},
				{
					Id:        "Event-1",
					MissionId: "Mission-1",
					Time:      2,
					Cancelled: false,
					Execute: func(m *Mission, gc *gamecomm.GameChannels) {
						msgChan <- eventMsg{
							msg: "Message 2",
						}
					},
				},
			},
			wants: []eventMsg{
				{
					msg:  "Message 2",
					time: 1,
				},
				{
					msg:  "Message 1",
					time: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameChannels := &gamecomm.GameChannels{
				WorldChannel:   make(chan gamecomm.WorldCommand),
				CorpChannel:    make(chan gamecomm.CorpCommand),
				MissionChannel: make(chan gamecomm.MissionCommand),
			}

			eventScheduler := &DefaultEventScheduler{
				events:       make(map[string]*Event),
				queue:        make(EventQueue, 0),
				missions:     make(map[string]*Mission),
				gameClock:    gameclock.NewGameClock(0, 1),
				gameChannels: gameChannels,
			}

			eventScheduler.idGenerator = func(e *Event) error {
				return nil
			}

			for _, e := range tt.events {
				_, err := eventScheduler.Schedule(e)
				if err != nil {
					t.Fatal(err)
				}
			}

			go eventScheduler.gameClock.StartTime()
			go eventScheduler.Run()

			for i := 0; i < len(tt.events); i++ {
				eMsg := <-msgChan

				currTime := eventScheduler.gameClock.GetCurrentTime()

				assert.Equal(t, eMsg.msg, tt.wants[i].msg)
				assert.Greater(t, currTime, tt.wants[i].time)
			}

		})
	}

	defer close(msgChan)

	msgChan = make(chan eventMsg)
	defer close(msgChan)

	// cancel event and it should not run
	t.Run("Cancelled Event Should Not Run", func(t *testing.T) {

		correctMsg := "Message 2"

		event1 := &Event{
			Id:        "Event-1",
			MissionId: "Mission-1",
			Time:      1,
			Cancelled: true,
			Execute: func(m *Mission, gc *gamecomm.GameChannels) {
				msgChan <- eventMsg{
					msg: "Message 1",
				}
			},
		}

		event2 := &Event{
			Id:        "Event-1",
			MissionId: "Mission-1",
			Time:      10,
			Cancelled: false,
			Execute: func(m *Mission, gc *gamecomm.GameChannels) {
				msgChan <- eventMsg{
					msg: correctMsg,
				}
			},
		}

		gameChannels := &gamecomm.GameChannels{
			WorldChannel:   make(chan gamecomm.WorldCommand),
			CorpChannel:    make(chan gamecomm.CorpCommand),
			MissionChannel: make(chan gamecomm.MissionCommand),
		}

		eventScheduler := &DefaultEventScheduler{
			events:       make(map[string]*Event),
			queue:        make(EventQueue, 0),
			missions:     make(map[string]*Mission),
			gameClock:    gameclock.NewGameClock(0, 100),
			gameChannels: gameChannels,
		}

		eventScheduler.idGenerator = func(e *Event) error {
			return nil
		}

		_, err := eventScheduler.Schedule(event1)
		if err != nil {
			t.Fatal(err)
		}
		_, err = eventScheduler.Schedule(event2)
		if err != nil {
			t.Fatal(err)
		}

		go eventScheduler.gameClock.StartTime()
		go eventScheduler.Run()

		eMsg := <-msgChan

		assert.Equal(t, eMsg.msg, correctMsg)
	})

	// high volume of events
	// idempotency

}

func BenchmarkRun(b *testing.B) {

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var wg sync.WaitGroup
		numEvents := 1000

		eventScheduler := &DefaultEventScheduler{
			events:    make(map[string]*Event),
			queue:     make(EventQueue, 0),
			missions:  make(map[string]*Mission),
			gameClock: gameclock.NewGameClock(0, 2),
		}

		eventScheduler.idGenerator = func(e *Event) error {
			return nil
		}

		var events []*Event

		b.StartTimer()
		for i := 0; i < numEvents; i++ {
			event := &Event{
				Id:        fmt.Sprintf("Event-%v", i),
				MissionId: "Mission-1",
				Time:      gameclock.GameTime((i + 1) * 2),
				Cancelled: false,
				Execute: func(m *Mission, gc *gamecomm.GameChannels) {
				},
			}

			events = append(events, event)
		}

		for _, event := range events {
			wg.Add(1)
			go func(e *Event, eventScheduler *DefaultEventScheduler) {
				defer wg.Done()

				_, err := eventScheduler.Schedule(e)
				if err != nil {
					b.Error(err)
				}

			}(event, eventScheduler)
		}

		go eventScheduler.gameClock.StartTime()
		go eventScheduler.Run()

		wg.Wait()

		b.StopTimer()

	}

}
