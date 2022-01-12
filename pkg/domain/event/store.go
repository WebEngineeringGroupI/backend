package event

import (
	"context"
	"sort"
)

// Stream represents an array of events
type Stream struct {
	events  []Event
	version int
}

func (e Stream) Events() []Event {
	return e.events
}

func (e Stream) Version() int {
	return e.version
}

// Len implements sort.Interface
func (e Stream) Len() int {
	return len(e.events)
}

// Swap implements sort.Interface
func (e Stream) Swap(i, j int) {
	e.events[i], e.events[j] = e.events[j], e.events[i]
}

// Less implements sort.Interface
func (e Stream) Less(i, j int) bool {
	return e.events[i].EventVersion() < e.events[j].EventVersion()
}

func (e *Stream) Append(events ...Event) {
	e.events = append(e.events, events...)
	sort.Sort(e)
	e.version = e.events[len(e.events)-1].EventVersion()
}

func StreamFrom(events []Event) *Stream {
	stream := &Stream{
		events: events,
	}
	sort.Sort(stream)
	stream.version = stream.events[stream.Len()-1].EventVersion()
	return stream
}

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
// Store provides an abstraction for the repository to save data
type Store interface {
	// Save the provided events to the store
	Append(ctx context.Context, identity string, events ...Event) error

	// Load the history of events.
	Load(ctx context.Context, identity string) (*Stream, error)
}
