package url

import (
	"context"
	"errors"
)

var ErrUnableToValidateURLs = errors.New("unable to validate URLs")

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type Validator interface {
	ValidateURLs(ctx context.Context, url []string) (bool, error)
}
