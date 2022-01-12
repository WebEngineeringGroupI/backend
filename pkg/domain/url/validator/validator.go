package validator

import (
	"context"
	"log"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Service struct {
	brokerReceiver redirector.ExternalBrokerReceiver
	brokerSender   redirector.ExternalBrokerSender
	urlValidator   url.Validator
	serializer     event.Serializer
	clock          event.Clock
}

func (s *Service) Start(ctx context.Context) error {
	eventsDataCh, err := s.brokerReceiver.ReceiveEvents(ctx)
	if err != nil {
		return err
	}

	for eventData := range eventsDataCh {
		evt, err := s.serializer.UnmarshalEvent(eventData)
		if err != nil {
			log.Printf("unable to retrieve event from data: %s", err)
			continue
		}
		s.handleEvent(ctx, evt)
	}
	return nil
}

func (s *Service) handleEvent(ctx context.Context, evt event.Event) {
	switch e := evt.(type) {
	case *url.ShortURLCreated:
		isValid, err := s.urlValidator.ValidateURLs(ctx, []string{e.OriginalURL})
		if err != nil {
			log.Printf("unable to validate URL %s: %s", e.OriginalURL, err)
			return
		}
		if isValid {
			log.Printf("validated url: %s", e.OriginalURL)
			s.sendEvent(ctx, &url.ShortURLVerified{
				Base: event.Base{
					ID:      e.EntityID(),
					Version: e.EventVersion() + 1,
					At:      s.clock.Now(),
				},
			})
		}
	case *url.LoadBalancedURLCreated:
		for idx, originalURL := range e.OriginalURLs {
			isValid, err := s.urlValidator.ValidateURLs(ctx, []string{originalURL})
			if err != nil {
				log.Printf("unable to validate URL %s: %s", originalURL, err)
				return
			}
			if isValid {
				log.Printf("validated url: %s", originalURL)
				s.sendEvent(ctx, &url.LoadBalancedURLVerified{
					Base: event.Base{
						ID:      e.EntityID(),
						Version: e.EventVersion() + idx + 1,
						At:      s.clock.Now(),
					},
					VerifiedURL: originalURL,
				})
			}
		}
	}
}

func (s *Service) sendEvent(ctx context.Context, event event.Event) {
	data, err := s.serializer.MarshalEvent(event)
	if err != nil {
		log.Printf("unable to marshal event to send it: %s", err)
		return
	}

	err = s.brokerSender.SendEvents(ctx, data)
	if err != nil {
		log.Printf("unable to send event: %s", err)
	}
}

func NewService(brokerReceiver redirector.ExternalBrokerReceiver, brokerSender redirector.ExternalBrokerSender, urlValidator url.Validator, serializer event.Serializer, clock event.Clock) *Service {
	return &Service{
		brokerReceiver: brokerReceiver,
		brokerSender:   brokerSender,
		urlValidator:   urlValidator,
		serializer:     serializer,
		clock:          clock,
	}
}
