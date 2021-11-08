package safebrowsing

import (
	"fmt"

	"github.com/google/safebrowsing"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Validator struct {
	safebrowser *safebrowsing.SafeBrowser
}

func (s *Validator) ValidateURL(aLongURL string) (bool, error) {
	return s.ValidateURLs([]string{aLongURL})
}

func (s *Validator) ValidateURLs(urls []string) (bool, error) {
	threats, err := s.safebrowser.LookupURLs(urls)
	if err != nil {
		return false, fmt.Errorf("%w: %s", url.ErrUnableToValidateURL, err)
	}

	for _, threat := range threats {
		if threat != nil {
			return false, nil
		}
	}

	return true, nil
}

func NewValidator(apiKey string) (*Validator, error) {
	safebrowser, err := safebrowsing.NewSafeBrowser(safebrowsing.Config{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &Validator{safebrowser}, nil
}
