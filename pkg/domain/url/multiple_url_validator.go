package url

import (
	"errors"
)

var ErrUnableToValidateURLs = errors.New("unable to validate multiple URLs")

type Validator interface {
	ValidateURLs(url []string) (bool, error)
}
