package postgres_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/serializer"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/serializer/json"
)

var _ = Describe("Infrastructure / Database / Postgres Outbox Source", func() {
	var (
		db              *postgres.DB
		eventSerializer serializer.Serializer
		ctx             context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		eventSerializer = json.NewSerializer(&Event1{}, &Event2{})
		var err error
		db, err = postgres.NewDB(connectionDetails(), eventSerializer)
		Expect(err).ToNot(HaveOccurred())
	})

	It("retrieves the events from the database", func() {
		identity := randomHash()
		err := db.Append(ctx, identity, Event1{
			Base: event.Base{
				ID:      identity,
				Version: 0,
				At:      time.Time{},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		events, err := db.PullEvents(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(events).To(ContainElement(
			WithTransform(
				eventFromPayloadWith(eventSerializer),
				Equal(&Event1{
					Base: event.Base{
						ID:      identity,
						Version: 0,
						At:      time.Time{},
					},
				}),
			),
		))
	})

	It("marks the events as sent in the database", func() {
		identity := randomHash()
		err := db.Append(ctx, identity, Event1{
			Base: event.Base{
				ID:      identity,
				Version: 0,
				At:      time.Time{},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		eventsBeforeSent, err := db.PullEvents(ctx)
		Expect(err).ToNot(HaveOccurred())

		err = db.MarkEventsAsSent(ctx, eventsBeforeSent)
		Expect(err).ToNot(HaveOccurred())

		eventsAfterSent, err := db.PullEvents(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(eventsAfterSent).ToNot(ContainElements(eventsBeforeSent))
	})
})

func eventFromPayloadWith(serializer serializer.Serializer) func(event *redirector.OutboxEvent) event.Event {
	return func(event *redirector.OutboxEvent) event.Event {
		unmarshalEvent, err := serializer.UnmarshalEvent(event.Payload)
		ExpectWithOffset(2, err).ToNot(HaveOccurred())
		return unmarshalEvent
	}
}
