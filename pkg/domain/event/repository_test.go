package event_test

import (
	"context"
	"log"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/eventstore/inmemory"
)

var _ = Describe("Domain / Event repository", func() {
	var (
		repository event.Repository
		ctx        context.Context
		ctrl       *gomock.Controller
		broker     *mocks.MockBroker
	)
	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		broker = mocks.NewMockBroker(ctrl)
		repository = event.NewRepository(&SomeEntity{}, inmemory.NewEventStore(), broker)
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	It("is able to save the events in the repository", func() {
		event := &SomeEntityCreated{Base: event.Base{ID: "1"}}
		broker.EXPECT().Publish(event)

		err := repository.Save(ctx, event)

		Expect(err).ToNot(HaveOccurred())
	})

	It("is able to retrieve the entity in the final state with all the events applied", func() {
		event := &SomeEntityCreated{Base: event.Base{ID: "1", Version: 2}}
		broker.EXPECT().Publish(event)

		err := repository.Save(ctx, event)
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
