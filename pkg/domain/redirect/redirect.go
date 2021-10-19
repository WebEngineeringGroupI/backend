package redirect

import (
	`errors`
	`fmt`
	"strings"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Redirector struct {
	repository url.ShortURLRepository
}

func NewRedirector(repository url.ShortURLRepository) *Redirector {
	return &Redirector{repository: repository}
}

func (r *Redirector) ReturnOriginalURL(hash string) (string, error) {
	//FindByHash
	shortURL, err := r.repository.FindByHash(hash)
	if errors.Is(err, url.ErrShortURLNotFound) {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("unexpected error retrieving original URL: %w", err)
	}

	aLongURL := shortURL.LongURL
	if !strings.HasPrefix(aLongURL, "http://") && !strings.HasPrefix(aLongURL, "https://") {
		return "", nil
	}
	return aLongURL, nil
}
