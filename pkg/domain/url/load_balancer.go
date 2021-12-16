package url

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoURLsSpecified     = errors.New("no URLs specified")
	ErrTooMuchMultipleURLs = errors.New("too much multiple URLs")
	ErrValidURLNotFound    = errors.New("valid URL not found")
)

const maxNumberOfURLsToLoadBalance = 10

type LoadBalancedURLsRepository interface {
	FindLoadBalancedURLByHash(hash string) (*LoadBalancedURL, error)
	SaveLoadBalancedURL(urls *LoadBalancedURL) error
}

type LoadBalancer struct {
	repository LoadBalancedURLsRepository
}

type LoadBalancedURL struct {
	Hash     string
	LongURLs []OriginalURL
}

func (b *LoadBalancer) ShortURLs(urls []string) (*LoadBalancedURL, error) {
	if len(urls) == 0 {
		return nil, ErrNoURLsSpecified
	}
	if len(urls) > maxNumberOfURLsToLoadBalance {
		return nil, ErrTooMuchMultipleURLs
	}

	multipleShortURLs := &LoadBalancedURL{
		Hash:     hashFromURLs(urls),
		LongURLs: originalURLsFromRaw(urls),
	}

	err := b.repository.SaveLoadBalancedURL(multipleShortURLs)
	if err != nil {
		return nil, fmt.Errorf("error saving load-balanced URLs into repository: %w", err)
	}

	return multipleShortURLs, nil
}

func originalURLsFromRaw(urls []string) []OriginalURL {
	originalURLs := make([]OriginalURL, 0, len(urls))
	for _, aURL := range urls {
		originalURLs = append(originalURLs, OriginalURL{
			URL:     aURL,
			IsValid: true,
		})
	}
	return originalURLs
}

func hashFromURLs(urls []string) string {
	return hashFromURL(strings.Join(urls, ""))
}

func NewLoadBalancer(repository LoadBalancedURLsRepository) *LoadBalancer {
	return &LoadBalancer{
		repository: repository,
	}
}
