package redirector_test

import (
	"context"
	"errors"
	"log"
	"reflect"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector"
)

var _ = Describe("Domain / Event / Redirector", func() {
	var (
		outboxSource      *FakeOutboxSource
		externalBroker    *FakeExternalBroker
		ctx               context.Context
		cancel            context.CancelFunc
		redirectorService *redirector.Redirector
	)

	BeforeEach(func() {
		log.Default().SetOutput(GinkgoWriter)
		ctx, cancel = context.WithCancel(context.Background())
		outboxSource = &FakeOutboxSource{}
		externalBroker = &FakeExternalBroker{receivedEvents: [][]byte{}}
		redirectorService = redirector.NewRedirector(outboxSource, externalBroker, 100*time.Millisecond)
	})

	AfterEach(func() {
		cancel()
	})

	It("retrieves the events from the outbox, and sends them to the external broker", func() {
		outboxSource.shouldReturnEvents(&redirector.OutboxEvent{ID: 0, Payload: []byte("somePayload")})

		go redirectorService.Start(ctx)

		Eventually(func() [][]byte { return externalBroker.ReceivedEvents() }).Should(ContainElement([]byte("somePayload")))
	})

	When("it is unable to send the event to the external broker", func() {
		It("doesn't delete it", func() {
			outboxSource.shouldReturnEvents(
				&redirector.OutboxEvent{ID: 0, Payload: []byte("somePayload")},
				&redirector.OutboxEvent{ID: 1, Payload: []byte("somePayload2")},
			)
			externalBroker.shouldFailWhenSendingPayload([]byte("somePayload2"))

			go redirectorService.Start(ctx)

			Eventually(func() [][]byte { return externalBroker.ReceivedEvents() }).Should(ContainElement([]byte("somePayload")))
			Consistently(func() [][]byte { return externalBroker.ReceivedEvents() }).ShouldNot(ContainElement([]byte("somePayload2")))
			Expect(outboxSource.EventsToReturn()).To(ContainElement(&redirector.OutboxEvent{ID: 1, Payload: []byte("somePayload2")}))
		})
	})
})

type FakeOutboxSource struct {
	eventsToReturn []*redirector.OutboxEvent
	mutex          sync.Mutex
}

func (f *FakeOutboxSource) shouldReturnEvents(events ...*redirector.OutboxEvent) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.eventsToReturn = events
}

func (f *FakeOutboxSource) PullEvents(context.Context) ([]*redirector.OutboxEvent, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.eventsToReturn, nil
}

func (f *FakeOutboxSource) MarkEventsAsSent(ctx context.Context, eventsToRemove []*redirector.OutboxEvent) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	var events []*redirector.OutboxEvent
	for _, event := range f.eventsToReturn {
		eventShouldBeRemoved := false
		for _, eventToRemove := range eventsToRemove {
			if reflect.DeepEqual(event, eventToRemove) {
				eventShouldBeRemoved = true
				break
			}
		}
		if !eventShouldBeRemoved {
			events = append(events, event)
		}
	}
	f.eventsToReturn = events
	return nil
}

func (f *FakeOutboxSource) EventsToReturn() []*redirector.OutboxEvent {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.eventsToReturn
}

type FakeExternalBroker struct {
	receivedEvents    [][]byte
	payloadToFailWith []byte
	mutex             sync.RWMutex
}

func (f *FakeExternalBroker) ReceiveEvents(context.Context) (<-chan []byte, error) {
	panic("implement me")
}

func (f *FakeExternalBroker) shouldFailWhenSendingPayload(payload []byte) {
	f.payloadToFailWith = payload
}

func (f *FakeExternalBroker) SendEvents(ctx context.Context, eventData ...[]byte) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	for _, data := range eventData {
		if reflect.DeepEqual(f.payloadToFailWith, data) {
			return errors.New("error")
		}
		f.receivedEvents = append(f.receivedEvents, data)
	}
	return nil
}

func (f *FakeExternalBroker) ReceivedEvents() [][]byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.receivedEvents
}
