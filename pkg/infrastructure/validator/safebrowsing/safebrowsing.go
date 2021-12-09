package safebrowsing

import (
	"context"
	"fmt"

	"github.com/google/safebrowsing"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Validator struct {
	safebrowser *safebrowsing.SafeBrowser
}

func (s *Validator) ValidateURLs(ctx context.Context, urls []string) (bool, error) {
	threats, err := s.safebrowser.LookupURLsContext(ctx, urls)
	if err != nil {
		return false, fmt.Errorf("%w: %s", url.ErrUnableToValidateURLs, err)
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
