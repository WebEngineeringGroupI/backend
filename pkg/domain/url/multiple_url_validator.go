package url

import (
	"errors"
)

var ErrUnableToValidateURLs = errors.New("unable to validate URLs")

type Validator interface {
	ValidateURLs(url []string) (bool, error)
}
