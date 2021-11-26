package domain_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain"
)

var _ = Describe("Domain / Broker", func() {
	var (
		broker         *domain.Broker
		fakeSubscriber *FakeSubscriber
	)
	BeforeEach(func() {
		broker = domain.NewBroker()
		fakeSubscriber = &FakeSubscriber{}
	})

	When("a subscriber subscribes for events", func() {
		It("successfully subscribes for fake events", func() {
			broker.Subscribe(fakeSubscriber, &FakeEvent1{}, &FakeEvent2{})
			broker.Publish(&FakeEvent1{})
			broker.Publish(&FakeEvent2{})

			Eventually(fakeSubscriber.EventsHandled).Should(HaveLen(2))
			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(BeAssignableToTypeOf(&FakeEvent1{}), BeAssignableToTypeOf(&FakeEvent2{})))
		})
	})
	When("a subscriber is not subscribed for an event", func() {
		It("doesn't receive this event", func() {
			broker.Subscribe(fakeSubscriber, &FakeEvent1{})
			broker.Publish(&FakeEvent1{})
			broker.Publish(&FakeEvent2{})

			Eventually(fakeSubscriber.EventsHandled).Should(HaveLen(1))
			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(BeAssignableToTypeOf(&FakeEvent1{})))
		})
	})
	When("a subscriber doesn't specify the type of events it subscribes to", func() {
		It("receives all events", func() {
			broker.Subscribe(fakeSubscriber)
			broker.Publish(&FakeEvent1{})
			broker.Publish(&FakeEvent1{})
			broker.Publish(&FakeEvent2{})
			broker.Publish(&FakeEvent3{})

			Eventually(fakeSubscriber.EventsHandled).Should(HaveLen(4))
			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(
				BeAssignableToTypeOf(&FakeEvent1{}),
				BeAssignableToTypeOf(&FakeEvent1{}),
				BeAssignableToTypeOf(&FakeEvent2{}),
				BeAssignableToTypeOf(&FakeEvent3{}),
			))
		})
	})

	When("a subscriber unsubscribes to an event type", func() {
		It("doesn't receive events for this type", func() {
			broker.Subscribe(fakeSubscriber, &FakeEvent1{}, &FakeEvent2{})
			broker.Unsubscribe(fakeSubscriber, &FakeEvent1{})

			broker.Publish(&FakeEvent1{})
			broker.Publish(&FakeEvent2{})

			Eventually(fakeSubscriber.EventsHandled).Should(HaveLen(1))
			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(BeAssignableToTypeOf(&FakeEvent2{})))
		})
	})

	When("a subscriber subscribes to all events and then unsubscribes to all events", func() {
		It("doesn't receive any event", func() {
			broker.Subscribe(fakeSubscriber)
			broker.Unsubscribe(fakeSubscriber)

			broker.Publish(&FakeEvent1{})

			Consistently(fakeSubscriber.EventsHandled).Should(BeEmpty())
		})
	})

	When("a subscribes to some events and then unsubscribes to all events", func() {
		It("doesn't receive the events of this type", func() {
			broker.Subscribe(fakeSubscriber, &FakeEvent1{}, &FakeEvent2{})
			broker.Unsubscribe(fakeSubscriber)

			broker.Publish(&FakeEvent1{})
			broker.Publish(&FakeEvent2{})

			Consistently(fakeSubscriber.EventsHandled).Should(BeEmpty())
		})
	})
	When("a subscriber is subscribed twice for an event type", func() {
		It("doesn't receive duplicated events", func() {
			broker.Subscribe(fakeSubscriber, &FakeEvent1{}, &FakeEvent1{})

			broker.Publish(&FakeEvent1{})

			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(BeAssignableToTypeOf(&FakeEvent1{})))
		})
	})
	When("a subscriber is subscribed to all event types and then subscribes for a single event type", func() {
		It("doesn't receive duplicated events", func() {
			broker.Subscribe(fakeSubscriber)
			broker.Subscribe(fakeSubscriber, &FakeEvent1{})

			broker.Publish(&FakeEvent1{})

			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(BeAssignableToTypeOf(&FakeEvent1{})))
		})
	})
	When("a subscriber is subscribed to an event type, and then subscribes to all event types", func() {
		It("doesn't receive duplicated events", func() {
			broker.Subscribe(fakeSubscriber, &FakeEvent1{})
			broker.Subscribe(fakeSubscriber)

			broker.Publish(&FakeEvent1{})

			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(BeAssignableToTypeOf(&FakeEvent1{})))
		})
	})
	When("a subscriber is subscribed to all event types twice", func() {
		It("doesn't receive duplicated events", func() {
			broker.Subscribe(fakeSubscriber)
			broker.Subscribe(fakeSubscriber)

			broker.Publish(&FakeEvent1{})

			Eventually(fakeSubscriber.EventsHandled).Should(ConsistOf(BeAssignableToTypeOf(&FakeEvent1{})))
		})
	})
})

type FakeSubscriber struct {
	mutex         sync.RWMutex
	eventsHandled []domain.Event
}

func (f *FakeSubscriber) EventsHandled() []domain.Event {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.eventsHandled
}

func (f *FakeSubscriber) HandleEvent(event domain.Event) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.eventsHandled = append(f.eventsHandled, event)
}

type FakeEvent1 struct {
}

func (f *FakeEvent1) ID() string {
	return "id"
}

func (f *FakeEvent1) HappenedOn() time.Time {
	return time.Time{}
}

type FakeEvent2 struct {
}

func (f *FakeEvent2) ID() string {
	return "id"
}

func (f *FakeEvent2) HappenedOn() time.Time {
	return time.Time{}
}

type FakeEvent3 struct {
}

func (f *FakeEvent3) ID() string {
	return "id"
}

func (f *FakeEvent3) HappenedOn() time.Time {
	return time.Time{}
}
