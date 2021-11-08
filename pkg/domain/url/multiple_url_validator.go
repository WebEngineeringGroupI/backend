package url

import (
	"errors"
)

var ErrUnableToValidateURLs = errors.New("unable to validate multiple URLs")

type MultipleValidator interface {
	ValidateURLs(url []string) (bool, error)
}
