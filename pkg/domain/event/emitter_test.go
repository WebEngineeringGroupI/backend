package event_test

import (
	"context"
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
)

var _ = Describe("Domain / EventEmitter", func() {
	var (
		ctx     context.Context
		ctrl    *gomock.Controller
		clock   *mocks.MockClock
		uuid    *mocks.MockUUID
		outbox  *mocks.MockOutbox
		emitter event.Emitter
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		clock = mocks.NewMockClock(ctrl)
		uuid = mocks.NewMockUUID(ctrl)
		outbox = mocks.NewMockOutbox(ctrl)
		emitter = event.NewEmitter(outbox, clock, uuid)

		clock.EXPECT().Now().Return(time.Time{}).AnyTimes()
		uuid.EXPECT().NewUUID().Return("anUUID").AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("creates a ShortURLCreated event with the time and ID", func() {
		outbox.EXPECT().SaveEvent(ctx, event.NewShortURLCreated("anUUID", time.Time{}, "hash", "originalURL", false))

		err := emitter.EmitShortURLCreated(ctx, "hash", "originalURL", false)

		Expect(err).ToNot(HaveOccurred())
	})

	When("the outbox returns an error", func() {
		It("returns the error from the outbox", func() {
			outbox.EXPECT().
				SaveEvent(ctx, event.NewShortURLCreated("anUUID", time.Time{}, "hash", "originalURL", false)).
				Return(errors.New("unknown error"))

			err := emitter.EmitShortURLCreated(ctx, "hash", "originalURL", false)

			Expect(err).To(MatchError("unknown error"))
		})
	})
})
