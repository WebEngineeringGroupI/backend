package redirect

import (
	"errors"
	"fmt"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Redirector struct {
	repository url.ShortURLRepository
}

func (r *Redirector) ReturnOriginalURL(hash string) (string, error) {
	shortURL, err := r.repository.FindShortURLByHash(hash)
	if errors.Is(err, url.ErrShortURLNotFound) {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("unexpected error retrieving original URL: %w", err)
	}

	if !shortURL.OriginalURL.IsValid {
		return "", fmt.Errorf("the url '%s' is marked as invalid", shortURL.OriginalURL.URL)
	}

	return shortURL.OriginalURL.URL, nil
}

func NewRedirector(repository url.ShortURLRepository) *Redirector {
	return &Redirector{
		repository: repository,
	}
}
