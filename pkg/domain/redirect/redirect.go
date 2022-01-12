package redirect

import (
	"context"
	"errors"
	"fmt"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Redirector struct {
	repository event.Repository
	clock      event.Clock
}

func (r *Redirector) ReturnOriginalURL(ctx context.Context, hash string) (string, error) {
	shortURLEntity, version, err := r.repository.Load(ctx, hash)
	if errors.Is(err, event.ErrEntityNotFound) {
		return "", url.ErrShortURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("unexpected error retrieving original URL: %w", err)
	}

	shortURL, ok := shortURLEntity.(*url.ShortURL)
	if !ok {
		return "", event.ErrEntityNotFound
	}

	if !shortURL.OriginalURL.IsValid {
		return "", fmt.Errorf("the url '%s' is marked as invalid", shortURL.OriginalURL.URL)
	}

	err = r.repository.Save(ctx, &url.ShortURLClicked{
		Base: event.Base{
			ID:      shortURL.Hash,
			Version: version + 1,
			At:      r.clock.Now(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("error saving url clicked event in the repository")
	}

	return shortURL.OriginalURL.URL, nil
}

func NewRedirector(repository event.Repository, clock event.Clock) *Redirector {
	return &Redirector{
		repository: repository,
		clock:      clock,
	}
}
