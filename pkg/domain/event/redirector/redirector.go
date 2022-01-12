package redirector

import (
	"context"
	"log"
	"time"
)

type OutboxEvent struct {
	ID      int
	Payload []byte
}

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type OutboxSource interface {
	PullEvents(ctx context.Context) ([]*OutboxEvent, error)
	MarkEventsAsSent(ctx context.Context, events []*OutboxEvent) error
}

type ExternalBrokerSender interface {
	SendEvents(ctx context.Context, eventData ...[]byte) error
}

type ExternalBrokerReceiver interface {
	ReceiveEvents(ctx context.Context) (<-chan []byte, error)
}

type Redirector struct {
	outboxSource    OutboxSource
	externalBroker  ExternalBrokerSender
	pollingInterval time.Duration
}

func (r *Redirector) Start(ctx context.Context) {
	ticker := time.NewTicker(r.pollingInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.redirectEvents(ctx)
		}
	}
}

func (r *Redirector) redirectEvents(ctx context.Context) {
	events, err := r.outboxSource.PullEvents(ctx)
	if err != nil {
		log.Printf("error while pulling events: %s", err)
		return
	}

	var sentEvents []*OutboxEvent
	for _, event := range events {
		if err := r.externalBroker.SendEvents(ctx, event.Payload); err != nil {
			log.Printf("error sending event to external broker: %s", err)
			continue
		}
		sentEvents = append(sentEvents, event)
	}
	if err := r.outboxSource.MarkEventsAsSent(ctx, sentEvents); err != nil {
		log.Printf("error marking event as sent: %s", err)
	}
}

func NewRedirector(outboxSource OutboxSource, externalBroker ExternalBrokerSender, pollingInterval time.Duration) *Redirector {
	return &Redirector{
		outboxSource:    outboxSource,
		externalBroker:  externalBroker,
		pollingInterval: pollingInterval,
	}
}
