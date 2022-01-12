package event

import (
	"time"
)

// TODO(fede) Move this to another package

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type Clock interface {
	Now() time.Time
}
