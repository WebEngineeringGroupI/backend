package event

import (
	"context"
	"fmt"
	"log"
	"reflect"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
// Entity represents the aggregate root in the domain driven design sense.
// It represents the current state of the domain object and can be thought of
// as a left fold over events.
// The Entity should contain a version that should be updated for each new change
// applied by an event.
type Entity interface {
	// On will be called for each event; returns err if the event could not be
	// applied.
	On(event Event) error
}

type Repository interface {
	Save(ctx context.Context, events ...Event) error
	Load(ctx context.Context, entityID string) (Entity, int, error)
}

// repository provides the primary abstraction to saving and loading events from a specific aggregate
type repository struct {
	broker    Broker
	prototype reflect.Type
	store     Store
}

// New returns a new instance of the aggregate
func (r *repository) New() Entity {
	return reflect.New(r.prototype).Interface().(Entity)
}

// Save persists the events into the underlying Store
func (r *repository) Save(ctx context.Context, events ...Event) error {
	if len(events) == 0 {
		return nil
	}
	aggregateID := events[0].EntityID()
	err := r.store.Append(ctx, aggregateID, events...)

	for _, event := range events {
		r.broker.Publish(event)
	}
	return err
}

// Load retrieves the specified aggregate from the underlying store
func (r *repository) Load(ctx context.Context, entityID string) (Entity, int, error) {
	eventStream, err := r.store.Load(ctx, entityID)
	if err != nil {
		return nil, 0, err
	}

	eventCount := eventStream.Len()
	if eventCount == 0 {
		return nil, 0, fmt.Errorf("%w: unable to load entity [type=%v, id=%v]", ErrEntityNotFound, r.prototype.String(), entityID)
	}
	log.Printf("loaded %v event(s) for entity id, %v", eventCount, entityID)

	aggregate := r.New()
	for _, event := range eventStream.Events() {
		err = aggregate.On(event)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: aggregate was unable to handle event, %v", ErrUnhandledEvent, TypeOf(event))
		}
	}

	return aggregate, eventStream.Version(), nil
}

// NewRepository creates a new repository
func NewRepository(prototype Entity, store Store, broker Broker) Repository {
	eventType := reflect.TypeOf(prototype)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}

	return &repository{
		broker:    broker,
		prototype: eventType,
		store:     store,
	}
}
