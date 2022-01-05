package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
)

// EventStore provides an in-memory implementation of event.Store
type EventStore struct {
	mux        *sync.Mutex
	eventsByID map[string]*event.Stream
}

func (m *EventStore) Append(ctx context.Context, entityID string, records ...event.Event) error {
	if len(records) == 0 {
		return nil
	}

	if _, ok := m.eventsByID[entityID]; !ok {
		m.eventsByID[entityID] = event.StreamFrom(records)
		return nil
	}

	m.eventsByID[entityID].Append(records...)
	return nil
}

func (m *EventStore) Load(ctx context.Context, entityID string) (*event.Stream, error) {
	eventStream, ok := m.eventsByID[entityID]
	if !ok {
		return nil, fmt.Errorf("%w: no entity found with id, %v", event.ErrEntityNotFound, entityID)
	}
	return eventStream, nil
}

func NewEventStore() *EventStore {
	return &EventStore{
		mux:        &sync.Mutex{},
		eventsByID: map[string]*event.Stream{},
	}
}
