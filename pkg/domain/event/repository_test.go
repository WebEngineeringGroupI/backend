package event_test

import (
	"context"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/eventstore/inmemory"
)

var _ = Describe("Domain / Event repository", func() {
	var (
		repository event.Repository
		ctx        context.Context
	)
	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		repository = event.NewRepository(&SomeEntity{}, inmemory.NewEventStore())
		ctx = context.Background()
	})

	It("is able to save the events in the repository", func() {
		err := repository.Save(ctx, &SomeEntityCreated{Base: event.Base{
			ID: "1",
		}})
		Expect(err).ToNot(HaveOccurred())
	})

	It("is able to retrieve the entity in the final state with all the events applied", func() {
		err := repository.Save(ctx, &SomeEntityCreated{Base: event.Base{
			ID:      "1",
			Version: 2,
		}})
		Expect(err).ToNot(HaveOccurred())

		aggregate, version, err := repository.Load(ctx, "1")
		Expect(err).ToNot(HaveOccurred())
		Expect(version).To(Equal(2))
		Expect(aggregate).To(Equal(&SomeEntity{ID: "1", version: 2}))
	})
})

type SomeEntity struct {
	ID      string
	version int
}

func (s *SomeEntity) On(event event.Event) error {
	switch e := event.(type) {
	case *SomeEntityCreated:
		s.ID = e.ID
		s.version = e.EventVersion()
	}
	return nil
}

type SomeEntityCreated struct {
	event.Base
}
