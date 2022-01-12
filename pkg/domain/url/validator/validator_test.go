package validator_test

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
	urlmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/validator"
)

var _ = Describe("Validator", func() {
	var (
		validatorService       *validator.Service
		ctrl                   *gomock.Controller
		externalBrokerReceiver *mocks.MockExternalBrokerReceiver
		externalBrokerSender   *mocks.MockExternalBrokerSender
		urlValidator           *urlmocks.MockValidator
		clock                  *eventmocks.MockClock
		serializer             event.Serializer
		ctx                    context.Context
		logger                 *strings.Builder
	)
	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		externalBrokerReceiver = mocks.NewMockExternalBrokerReceiver(ctrl)
		externalBrokerSender = mocks.NewMockExternalBrokerSender(ctrl)
		urlValidator = urlmocks.NewMockValidator(ctrl)
		clock = eventmocks.NewMockClock(ctrl)
		serializer = json.NewSerializer(&url.ShortURLCreated{}, &url.LoadBalancedURLCreated{})
		logger = &strings.Builder{}
		log.Default().SetOutput(logger)

		validatorService = validator.NewService(externalBrokerReceiver, externalBrokerSender, urlValidator, serializer, clock)

		clock.EXPECT().Now().Return(time.Time{}).AnyTimes()
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	DescribeTable("retrieves events and validates them or not",
		func(eventToReceive event.Event, isValidURL bool, eventsToRespond ...event.Event) {
			externalBrokerReceiver.EXPECT().ReceiveEvents(ctx).Return(channelWithEvents(eventToReceive), nil)
			urlValidator.EXPECT().ValidateURLs(ctx, gomock.Any()).Return(isValidURL, nil).AnyTimes()
			for _, evt := range eventsToRespond {
				externalBrokerSender.EXPECT().SendEvents(ctx, [][]byte{eventPayload(evt)}).Return(nil)
			}

			err := validatorService.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			Consistently(logger.String()).ShouldNot(ContainSubstring("unable"))
		},
		Entry("retrieves a shortURLCreated event and is valid",
			shortURLCreatedEvent("someURL"), true, shortURLVerifiedEvent()),
		Entry("retrieves a shortURLCreated event and is not valid",
			shortURLCreatedEvent("someURL"), false),
		Entry("retrieves a loadBalancedURLCreated event and is valid",
			loadBalancedURLCreatedEvent([]string{"someURL1", "someURL2"}), true, loadBalancedURLVerifiedEvent("someURL1", 1), loadBalancedURLVerifiedEvent("someURL2", 2)),
		Entry("retrieves a loadBalancedURLCreated event and is not valid",
			loadBalancedURLCreatedEvent([]string{"someURL1", "someURL2"}), false),
	)
})

func loadBalancedURLVerifiedEvent(verifiedURL string, version int) *url.LoadBalancedURLVerified {
	return &url.LoadBalancedURLVerified{
		Base: event.Base{
			ID:      "someID",
			Version: version,
			At:      time.Time{},
		},
		VerifiedURL: verifiedURL,
	}
}

func loadBalancedURLCreatedEvent(originalURLs []string) *url.LoadBalancedURLCreated {
	return &url.LoadBalancedURLCreated{
		Base: event.Base{
			ID:      "someID",
			Version: 0,
			At:      time.Time{},
		},
		OriginalURLs: originalURLs,
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

func shortURLCreatedEvent(originalURL string) *url.ShortURLCreated {
	return &url.ShortURLCreated{
		Base: event.Base{
			ID:      "someID",
			Version: 0,
			At:      time.Time{},
		},
		OriginalURL: originalURL,
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
