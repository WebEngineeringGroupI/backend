package postgres_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/serializer/json"
)

var _ = Describe("Infrastructure / Database / Postgres Event Sourcing", func() {
	var (
		db  *postgres.DB
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = postgres.NewDB(connectionDetails(), json.NewSerializer(&Event1{}, &Event2{}))
		Expect(err).ToNot(HaveOccurred())
	})

	It("retrieves the events correctly from the database", func() {
		entityID := randomHash()
		err := db.Append(ctx, entityID,
			&Event1{
				Base: event.Base{
					ID:      entityID,
					Version: 0,
					At:      time.Time{},
				},
			},
			&Event2{
				Base: event.Base{
					ID:      entityID,
					Version: 1,
					At:      time.Time{},
				},
			},
		)
		Expect(err).ToNot(HaveOccurred())

		stream, err := db.Load(ctx, entityID)
		Expect(err).ToNot(HaveOccurred())
		Expect(stream.Events()).To(ConsistOf(
			&Event1{
				Base: event.Base{
					ID:      entityID,
					Version: 0,
					At:      time.Time{},
				},
			},
			&Event2{
				Base: event.Base{
					ID:      entityID,
					Version: 1,
					At:      time.Time{},
				},
			}),
		)
		Expect(stream.Version()).To(Equal(1))
	})

	When("saving a duplicated version of an entity", func() {
		It("doesn't save it, but saves the other events", func() {
			entityID := randomHash()
			err := db.Append(ctx, entityID,
				&Event1{
					Base: event.Base{
						ID:      entityID,
						Version: 0,
						At:      time.Time{},
					},
				},
				&Event2{
					Base: event.Base{
						ID:      entityID,
						Version: 1,
						At:      time.Time{},
					},
				},
				&Event1{
					Base: event.Base{
						ID:      entityID,
						Version: 1,
						At:      time.Time{},
					},
				},
			)

			Expect(err).To(MatchError(ContainSubstring("unable to insert event in database, check the version of the events")))

			stream, err := db.Load(ctx, entityID)

			Expect(err).To(MatchError(event.ErrEntityNotFound))
			Expect(stream).To(BeNil())
		})
	})
})

type Event1 struct {
	event.Base
}
type Event2 struct {
	event.Base
}
