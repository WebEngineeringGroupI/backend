package validationsaver_test

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	eventmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/serializer/json"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/validationsaver"
)

var _ = Describe("ValidationSaver", func() {
	var (
		ctx                    context.Context
		ctrl                   *gomock.Controller
		validationSaverService *validationsaver.Service
		eventRepo              *eventmocks.MockRepository
		brokerReceiver         *mocks.MockExternalBrokerReceiver
		logger                 *strings.Builder
	)
	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		eventRepo = eventmocks.NewMockRepository(ctrl)
		brokerReceiver = mocks.NewMockExternalBrokerReceiver(ctrl)
		logger = &strings.Builder{}
		log.Default().SetOutput(logger)

		validationSaverService = validationsaver.NewService(eventRepo, brokerReceiver, json.NewSerializer(&url.ShortURLVerified{}, &url.LoadBalancedURLVerified{}))
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	DescribeTable("receives a validation event", func(eventReceived event.Event) {
		brokerReceiver.EXPECT().ReceiveEvents(ctx).Return(channelWithEvents(eventReceived), nil)
		eventRepo.EXPECT().Save(ctx, eventReceived).Return(nil)

		err := validationSaverService.Start(ctx)

		Expect(err).ToNot(HaveOccurred())
		Consistently(logger.String()).ShouldNot(ContainSubstring("unable"))
	},
		Entry("receives a shortURLVerified event", shortURLVerifiedEvent()),
		Entry("receives a loadBalancedURLVerified event", loadBalancedURLVerifiedEvent()),
	)

	When("receives a single url validated event", func() {
		It("saves the aggregate data in the database", func() {

		})
	})
})

func loadBalancedURLVerifiedEvent() *url.LoadBalancedURLVerified {
	return &url.LoadBalancedURLVerified{
		Base: event.Base{
			ID:      "someID",
			Version: 1,
			At:      time.Time{},
		},
		VerifiedURL: "someURL",
	}
}

func shortURLVerifiedEvent() *url.ShortURLVerified {
	return &url.ShortURLVerified{
		Base: event.Base{
			ID:      "someID",
			Version: 1,
			At:      time.Time{},
		},
	}
}

func eventPayload(event event.Event) []byte {
	data, _ := json.NewSerializer(event).MarshalEvent(event)
	return data
}

func channelWithEvents(events ...event.Event) chan []byte {
	recvChannel := make(chan []byte, len(events))
	go func() {
		defer close(recvChannel)
		for _, event := range events {
			recvChannel <- eventPayload(event)
		}
	}()
	return recvChannel
}
