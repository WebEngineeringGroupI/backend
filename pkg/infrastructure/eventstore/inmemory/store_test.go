package inmemory_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/eventstore/inmemory"
)

var _ = Describe("EventStore / InMemory Store", func() {
	var (
		store *inmemory.EventStore
		ctx   context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
		store = inmemory.NewEventStore()
	})

	It("appends the events correctly", func() {
		err := store.Append(context.Background(), "someID",
			&SomeEvent1{Base: event.Base{
				ID:      "someID",
				Version: 0,
				At:      time.Time{},
			}},
			&SomeEvent2{Base: event.Base{
				ID:      "someID",
				Version: 1,
				At:      time.Time{},
			}},
		)

		Expect(err).ToNot(HaveOccurred())
	})

	It("retrieves the events correctly", func() {
		err := store.Append(context.Background(), "someID",
			&SomeEvent1{Base: event.Base{
				ID:      "someID",
				Version: 0,
				At:      time.Time{},
			}},
			&SomeEvent2{Base: event.Base{
				ID:      "someID",
				Version: 1,
				At:      time.Time{},
			}},
		)

		Expect(err).ToNot(HaveOccurred())

		stream, err := store.Load(ctx, "someID")
		Expect(err).ToNot(HaveOccurred())
		Expect(stream.Version()).To(Equal(1))
		Expect(stream.Events()).To(ConsistOf(
			&SomeEvent1{Base: event.Base{
				ID:      "someID",
				Version: 0,
				At:      time.Time{},
			}},
			&SomeEvent2{Base: event.Base{
				ID:      "someID",
				Version: 1,
				At:      time.Time{},
			}},
		))
	})

})

type SomeEvent1 struct {
	event.Base
}

type SomeEvent2 struct {
	event.Base
}
