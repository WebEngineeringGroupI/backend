package event

import (
	"context"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type Outbox interface {
	SaveEvent(ctx context.Context, event Event) error
}

type Clock interface {
	Now() time.Time
}

type UUID interface {
	NewUUID() string
}

type Emitter interface {
	EmitShortURLCreated(ctx context.Context, hash string, originalURL string, isValid bool) error
}
