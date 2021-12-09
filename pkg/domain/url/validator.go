package url

import (
	"context"
	"errors"
)

var ErrUnableToValidateURLs = errors.New("unable to validate URLs")

type Validator interface {
	ValidateURLs(ctx context.Context, url []string) (bool, error)
}
