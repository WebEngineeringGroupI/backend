package event

import (
	"errors"
)

var (
	ErrEntityNotFound   = errors.New("entity not found")
	ErrUnhandledEvent   = errors.New("unhandled event")
	ErrUnableToEncode   = errors.New("unable to encode event")
	ErrUnableToDecode   = errors.New("unable to decode event")
	ErrUnknownEventType = errors.New("unknown event type")
)
