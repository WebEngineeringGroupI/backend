package url

import (
	"errors"
)

var ErrUnableToValidateURL = errors.New("unable to validate URL")

type Validator interface {
	ValidateURL(url string) (bool, error)
}
