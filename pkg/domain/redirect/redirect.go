package redirect

import (
	"errors"
	"fmt"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Redirector struct {
	repository url.ShortURLRepository
	validator  url.Validator
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

	isValidURL, err := r.validator.ValidateURL(shortURL.LongURL)
	if err != nil {
		return "", err
	}
	if !isValidURL {
		return "", fmt.Errorf("the url '%s' is marked as invalid", shortURL.LongURL)
	}

	return shortURL.LongURL, nil
}

func NewRedirector(repository url.ShortURLRepository, validator url.Validator) *Redirector {
	return &Redirector{
		repository: repository,
		validator:  validator,
	}
}
