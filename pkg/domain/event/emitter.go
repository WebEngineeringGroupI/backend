package event

import (
	"context"
)

type emitter struct {
	clock  Clock
	uuid   UUID
	outbox Outbox
}

func (e *emitter) EmitShortURLCreated(ctx context.Context, hash string, originalURL string) error {
	return e.outbox.SaveEvent(ctx, NewShortURLCreated(e.uuid.NewUUID(), e.clock.Now(), hash, originalURL))
}

func (e *emitter) EmitLoadBalancedURLCreated(ctx context.Context, hash string, originalURLs []string) error {
	return e.outbox.SaveEvent(ctx, NewLoadBalancedURLCreated(e.uuid.NewUUID(), e.clock.Now(), hash, originalURLs))
}

func NewEmitter(outbox Outbox, clock Clock, uuid UUID) Emitter {
	return &emitter{clock: clock, uuid: uuid, outbox: outbox}
}
