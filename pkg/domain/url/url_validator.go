package url

import (
	"errors"
	"strings"
)

var ErrUnableToValidateURL = errors.New("unable to validate URL")

type Validator interface {
	ValidateURL(url string) (bool, error)
}

type httpValidator struct{}

func (h *httpValidator) ValidateURL(aLongURL string) (bool, error) {
	isValidURL := strings.HasPrefix(aLongURL, "http://") || strings.HasPrefix(aLongURL, "https://")
	return isValidURL, nil
}

func NewHTTPValidator() Validator {
	return &httpValidator{}
}
