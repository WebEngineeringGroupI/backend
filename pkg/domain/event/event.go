package event

import (
	"time"
)

//Event represents something that has happened in the domain and can be sent
//through the Broker and handled by a Subscriber
type Event interface {
	// EntityID represents the ID of the entity referenced by the event.
	EntityID() string

	// EventVersion contains the version number of this event
	EventVersion() int

	// HappenedOn represents when the event occurred
	HappenedOn() time.Time
}

// Base provides a default implementation of an Event
type Base struct {
	// ID contains the EntityID
	ID string

	// Version contains the EventVersion
	Version int

	// At contains when it HappenedOn
	At time.Time
}

func (b Base) EntityID() string {
	return b.ID
}

func (b Base) EventVersion() int {
	return b.Version
}

func (b Base) HappenedOn() time.Time {
	return b.At
}
