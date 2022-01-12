package json

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
)

type jsonDomainEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Serializer provides a simple serializer implementation
type Serializer struct {
	eventTypes map[string]reflect.Type
}

// Bind registers the specified events with the serializer; may be called more than once
func (j *Serializer) Bind(events ...event.Event) {
	for _, anEvent := range events {
		typeOfEvent := reflect.TypeOf(anEvent)
		if typeOfEvent.Kind() == reflect.Ptr {
			typeOfEvent = typeOfEvent.Elem()
		}
		j.eventTypes[event.TypeOf(anEvent)] = typeOfEvent
	}
}

// MarshalEvent converts an event into its persistent type, Record
func (j *Serializer) MarshalEvent(v event.Event) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(jsonDomainEvent{
		Type: event.TypeOf(v),
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", event.ErrUnableToEncode, err)
	}

	return data, nil
}

// UnmarshalEvent converts the persistent type, Record, into an Event instance
func (j *Serializer) UnmarshalEvent(data []byte) (event.Event, error) {
	anEvent := jsonDomainEvent{}
	err := json.Unmarshal(data, &anEvent)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to unmarshal data: %s", event.ErrUnableToDecode, err)
	}

	eventType, ok := j.eventTypes[anEvent.Type]
	if !ok {
		return nil, fmt.Errorf("%w: %v", event.ErrUnknownEventType, anEvent.Type)
	}

	eventInstance := reflect.New(eventType).Interface()
	err = json.Unmarshal(anEvent.Data, eventInstance)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to unmarshal event data: %s", event.ErrUnableToDecode, err)
	}

	return eventInstance.(event.Event), nil
}

// NewSerializer constructs a new Serializer and populates it with the specified events.
// Bind may be subsequently called to add more events.
func NewSerializer(events ...event.Event) *Serializer {
	serializer := &Serializer{
		eventTypes: map[string]reflect.Type{},
	}
	serializer.Bind(events...)

	return serializer
}
