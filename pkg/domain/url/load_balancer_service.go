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

type LoadBalancerService struct {
	repository event.Repository
	clock      event.Clock
}

type LoadBalancedURL struct {
	Hash     string
	LongURLs []OriginalURL
}

func (l *LoadBalancedURL) On(evt event.Event) error {
	switch e := evt.(type) {
	case *LoadBalancedURLCreated:
		l.Hash = e.EntityID()
		l.LongURLs = []OriginalURL{}
		for _, url := range e.OriginalURLs {
			l.LongURLs = append(l.LongURLs, OriginalURL{
				IsValid: false,
				URL:     url,
			})
		}
	case *LoadBalancedURLVerified:
		l.LongURLs = verifyLongURLFromList(l.LongURLs, e.VerifiedURL)
	default:
		return event.ErrUnhandledEvent
	}
	return nil
}

func verifyLongURLFromList(original []OriginalURL, verifiedURL string) []OriginalURL {
	newList := make([]OriginalURL, 0, len(original))
	for _, url := range original {
		if url.URL != verifiedURL {
			newList = append(newList, url)
			continue
		}
		newList = append(newList, OriginalURL{
			URL:     url.URL,
			IsValid: true,
		})
	}
	return newList
}

func (b *LoadBalancerService) ShortURLs(ctx context.Context, urls []string) (*LoadBalancedURL, error) {
	if len(urls) == 0 {
		return nil, ErrNoURLsSpecified
	}
	if len(urls) > maxNumberOfURLsToLoadBalance {
		return nil, ErrTooMuchMultipleURLs
	}

	hash := hashFromURLs(urls)
	entity, _, err := b.repository.Load(ctx, hash)
	if err == nil {
		loadBalancedURL, ok := entity.(*LoadBalancedURL)
		if !ok {
			return nil, fmt.Errorf("unknown entity type loaded while load balancing urls: %w", err)
		}
		return loadBalancedURL, nil
	}

	events := []event.Event{
		&LoadBalancedURLCreated{
			Base: event.Base{
				ID:      hash,
				Version: 0,
				At:      b.clock.Now(),
			},
			OriginalURLs: urls,
		},
	}

	err = b.repository.Save(ctx, events...)
	if err != nil {
		return nil, fmt.Errorf("error saving load-balanced URLs into repository: %w", err)
	}

	url := &LoadBalancedURL{}
	for _, e := range events {
		err := url.On(e)
		if err != nil {
			return nil, err
		}
	}

	return url, nil
}

func hashFromURLs(urls []string) string {
	return hashFromURL(strings.Join(urls, ""))
}

func NewLoadBalancer(repository event.Repository, clock event.Clock) *LoadBalancerService {
	return &LoadBalancerService{
		repository: repository,
		clock:      clock,
	}
}
