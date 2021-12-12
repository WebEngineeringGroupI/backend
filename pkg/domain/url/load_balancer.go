package url

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
)

var (
	ErrNoURLsSpecified     = errors.New("no URLs specified")
	ErrTooMuchMultipleURLs = errors.New("too much multiple URLs")
	ErrValidURLNotFound    = errors.New("valid URL not found")
)

const maxNumberOfURLsToLoadBalance = 10

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type LoadBalancedURLsRepository interface {
	FindLoadBalancedURLByHash(ctx context.Context, hash string) (*LoadBalancedURL, error)
	SaveLoadBalancedURL(ctx context.Context, urls *LoadBalancedURL) error
}

type LoadBalancer struct {
	repository LoadBalancedURLsRepository
	emitter    event.Emitter
}

type LoadBalancedURL struct {
	Hash     string
	LongURLs []OriginalURL
}

func (b *LoadBalancer) ShortURLs(ctx context.Context, urls []string) (*LoadBalancedURL, error) {
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

	err := b.repository.SaveLoadBalancedURL(ctx, multipleShortURLs)
	if err != nil {
		return nil, fmt.Errorf("error saving load-balanced URLs into repository: %w", err)
	}
	err = b.emitter.EmitLoadBalancedURLCreated(ctx, multipleShortURLs.Hash, urls)
	if err != nil {
		return nil, fmt.Errorf("error emitting event of load balanced URLs created: %w", err)
	}

	return multipleShortURLs, nil
}

func originalURLsFromRaw(urls []string) []OriginalURL {
	originalURLs := make([]OriginalURL, 0, len(urls))
	for _, aURL := range urls {
		originalURLs = append(originalURLs, OriginalURL{
			URL:     aURL,
			IsValid: false,
		})
	}
	return originalURLs
}

func hashFromURLs(urls []string) string {
	return hashFromURL(strings.Join(urls, ""))
}

func NewLoadBalancer(repository LoadBalancedURLsRepository, emitter event.Emitter) *LoadBalancer {
	return &LoadBalancer{
		repository: repository,
		emitter:    emitter,
	}
}
