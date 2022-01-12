package validationsaver

import (
	"context"
	"fmt"
	"log"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Service struct {
	eventRepo      event.Repository
	brokerReceiver redirector.ExternalBrokerReceiver
	// FIXME(fede): Some refactor in the future, the serializer could be embedded in the broker receiver, thus, only receiving events and not byte slices.
	serializer event.Serializer
}

func (s *Service) Start(ctx context.Context) error {
	events, err := s.brokerReceiver.ReceiveEvents(ctx)
	if err != nil {
		return fmt.Errorf("unable to receive events from broker: %w", err)
	}
	for eventPayload := range events {
		evt, err := s.serializer.UnmarshalEvent(eventPayload)
		if err != nil {
			log.Printf("unable to unmarshal event from broker: %s", err)
		}
		s.handleEvent(ctx, evt)
	}
	return nil
}

func (s *Service) handleEvent(ctx context.Context, evt event.Event) {
	switch e := evt.(type) {
	case *url.ShortURLVerified, *url.LoadBalancedURLVerified:
		err := s.eventRepo.Save(ctx, e)
		if err != nil {
			log.Printf("unable to save event in the repository: %s", err)
		}
	}
}

func NewService(eventRepo event.Repository, brokerReceiver redirector.ExternalBrokerReceiver, eventSerializer event.Serializer) *Service {
	return &Service{
		eventRepo:      eventRepo,
		brokerReceiver: brokerReceiver,
		serializer:     eventSerializer,
	}
}
