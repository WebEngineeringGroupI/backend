package postgres_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"xorm.io/xorm"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/model"
)

var _ = Describe("Infrastructure / Database / Postgres / EventOutbox", func() {
	var (
		repository *postgres.DB
		session    *postgres.DBSession
		ctx        context.Context
		db         *xorm.Engine
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		repository, err = postgres.NewDB(connectionDetails())
		Expect(err).ToNot(HaveOccurred())
		db = postgresDB(connectionDetails())
		session = repository.Session()
	})

	AfterEach(func() {
		session.Close()
	})

	It("saves the event in the database", func() {
		err := session.SaveEvent(ctx, event.NewShortURLCreated("event_id", time.Time{}, "hash", "originalURL", false))
		Expect(err).ToNot(HaveOccurred())

		domainEvent := &model.DomainEvent{ID: "event_id"}
		has, err := db.Get(domainEvent)
		Expect(err).ToNot(HaveOccurred())
		Expect(has).To(BeTrue())
		Expect(domainEvent.ID).To(Equal("event_id"))
		Expect(domainEvent.Payload).To(Equal([]byte(`{"Hash": "hash", "EventID": "event_id", "IsValid": false, "Creation": "0001-01-01T00:00:00Z", "OriginalURL": "originalURL"}`)))
	})
})

func postgresDB(details postgres.ConnectionDetails) *xorm.Engine {
	engine, err := xorm.NewEngine("postgres", details.ConnectionString())
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return engine
}
