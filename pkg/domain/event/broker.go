package event

import (
	"reflect"
	"sync"
)

//Subscriber is a component that's able to receive an Event and will handle it.
//A Subscriber can be subscribed to the Broker through Broker.Subscribe to
//multiple event types, or all of them.
type Subscriber interface {
	HandleEvent(event Event)
}

//Broker has the behavior of a Message Broker (https://en.wikipedia.org/wiki/Message_broker)
//for a publish-subscribe pattern (https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern)
//where multiple subscribers can subscribe to all events, or single events and will only receive the
//desired events and act accordingly to them.
//
//Subscribers can unsubscribe to all events or single events, but there's no support for a subscriber that
//subscribes to all events, and then unsubscribes to only some of them.
type Broker struct {
	mutex              sync.RWMutex
	eventSubscriberMap map[string][]Subscriber
}

const allEventsID = "all_events"

//Publish publishes an event to all subscribers that are subscribed for this event type, or all event types.
func (b *Broker) Publish(event Event) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if subscribers, ok := b.eventSubscriberMap[TypeOf(event)]; ok {
		for _, subscriber := range subscribers {
			go subscriber.HandleEvent(event)
		}
	}
	for _, subscriber := range b.eventSubscriberMap[allEventsID] {
		go subscriber.HandleEvent(event)
	}
}

// TypeOf is a helper func that extracts the event type of the event along with the reflect. TypeOf of the event.
func TypeOf(event Event) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

//Subscribe subscribes a Subscriber for the specified event types passed as parameter.
//If no events types are specified, the subscriber is subscribed to all event types;
//keep in mind that once a subscriber is subscribed to all event types, it cannot be unsubscribed to
//a single event type.
//Duplicated subscriptions to the same event type will be ignored.
func (b *Broker) Subscribe(subscriber Subscriber, eventsToSubscribe ...Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(eventsToSubscribe) == 0 {
		b.unsubscribeFromAllEvents(subscriber)
		b.eventSubscriberMap[allEventsID] = append(b.eventSubscriberMap[allEventsID], subscriber)
	}

	for _, event := range eventsToSubscribe {
		eventType := TypeOf(event)
		if b.isSubscriberAlreadySubscribedToEventType(subscriber, eventType) {
			continue
		}
		b.eventSubscriberMap[eventType] = append(b.eventSubscriberMap[eventType], subscriber)
	}
}

//Unsubscribe will unsubscribe a Subscriber from the specified event types,
//or all of them if no event type is specified.
//Keep in mind that a Subscriber that's subscribed to all event types, cannot be
//unsubscribed to a single event types.
func (b *Broker) Unsubscribe(subscriberToUnsubscribe Subscriber, events ...Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(events) == 0 {
		b.unsubscribeFromAllEvents(subscriberToUnsubscribe)
	}

	for _, event := range events {
		b.unsubscribeForSingleEventType(subscriberToUnsubscribe, TypeOf(event))
	}
}

func (b *Broker) isSubscriberAlreadySubscribedToEventType(subscriber Subscriber, eventType string) bool {
	for _, subscriberForAllEventTypes := range b.eventSubscriberMap[allEventsID] {
		if subscriberForAllEventTypes == subscriber {
			return true
		}
	}

	subscribersForEventType, ok := b.eventSubscriberMap[eventType]
	if !ok {
		return false
	}

	for _, subscriberForEventType := range subscribersForEventType {
		if subscriberForEventType == subscriber {
			return true
		}
	}
	return false
}

func (b *Broker) unsubscribeFromAllEvents(subscriberToUnsubscribe Subscriber) {
	for eventType, subscribersForAType := range b.eventSubscriberMap {
		newSubscribersListForAType := make([]Subscriber, 0, len(subscribersForAType))
		for _, singleSubscriberForAType := range subscribersForAType {
			if singleSubscriberForAType != subscriberToUnsubscribe {
				newSubscribersListForAType = append(newSubscribersListForAType, singleSubscriberForAType)
			}
		}
		b.eventSubscriberMap[eventType] = newSubscribersListForAType
	}
}

func (b *Broker) unsubscribeForSingleEventType(subscriberToUnsubscribe Subscriber, eventType string) {
	if subscribers, ok := b.eventSubscriberMap[eventType]; ok {
		newSubscribersList := make([]Subscriber, 0, len(subscribers))
		for _, subscriberInList := range subscribers {
			if subscriberInList != subscriberToUnsubscribe {
				newSubscribersList = append(newSubscribersList, subscriberInList)
			}
		}
		b.eventSubscriberMap[eventType] = newSubscribersList
	}
}

//NewBroker creates a new broker that will handle subscriptions and event sending.
func NewBroker() *Broker {
	return &Broker{
		eventSubscriberMap: map[string][]Subscriber{
			allEventsID: {},
		},
	}
}
