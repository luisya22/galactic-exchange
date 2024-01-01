package mission

import (
	"container/heap"

	"github.com/luisya22/galactic-exchange/gameclock"
)

type EventQueue []*Event

func (eq EventQueue) Len() int {
	return len(eq)
}

func (eq EventQueue) Less(i, j int) bool {
	return eq[i].Time.Before(eq[j].Time)
}

func (eq EventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].Index = i
	eq[j].Index = j
}

func (eq *EventQueue) Push(e any) {
	n := len(*eq)
	event := e.(*Event)
	event.Index = n
	*eq = append(*eq, event)
}

func (eq *EventQueue) Pop() any {
	old := *eq
	n := len(old)
	event := old[n-1]
	event.Index = -1
	*eq = old[0 : n-1]
	return event
}

func (eq *EventQueue) Update(event *Event, t gameclock.GameTime, cancelled bool) {
	event.Time = t
	event.Cancelled = cancelled
	heap.Fix(eq, event.Index)
}
