package redirect

import (
	"errors"
	"fmt"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"math/rand"
)

var (
	ValidURLNotFound = errors.New("valid URL not found")
)

type MultipleShortURLsRepository interface {
	FindOriginalURLsForHash(hash string) ([]url.OriginalURL, error)
}

type MultipleURLRedirector struct {
	repository MultipleShortURLsRepository
}

func (r *MultipleURLRedirector) ReturnValidOriginalURL(hash string) (string, error) {
	originalURLs, err := r.repository.FindOriginalURLsForHash(hash)
	if err != nil {
		return "", err
	}

	validURLs := r.filterValidURLs(originalURLs)
	if len(validURLs) == 0 {
		return "", fmt.Errorf("there are no valid URLs to redirect to")
	}

	randomIndex := rand.Intn(len(validURLs))
	return validURLs[randomIndex], nil
}

func (r *MultipleURLRedirector) filterValidURLs(originalURLs []url.OriginalURL) []string {
	validURLs := []string{}
	for _, aURL := range originalURLs {
		if aURL.IsValid {
			validURLs = append(validURLs, aURL.URL)
		}
	}
	return validURLs
}

func NewMultipleURLRedirector(repository MultipleShortURLsRepository) *MultipleURLRedirector {
	return &MultipleURLRedirector{repository: repository}
}
