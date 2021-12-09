package domain

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type EventOutbox interface {
	SaveEvent(ctx context.Context, event Event) error
}
