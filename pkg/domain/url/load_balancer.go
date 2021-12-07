package url

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoURLsSpecified     = errors.New("no URLs specified")
	ErrTooMuchMultipleURLs = errors.New("too much multiple URLs")
)

const maxNumberOfURLsToLoadBalance = 10

type MultipleShortURLsRepository interface {
	Save(urls *MultipleShortURLs) error
}

type LoadBalancer struct {
	repository MultipleShortURLsRepository
}

type MultipleShortURLs struct {
	Hash     string
	LongURLs []string
}

func (b *LoadBalancer) ShortURLs(urls []string) (*MultipleShortURLs, error) {
	if len(urls) == 0 {
		return nil, ErrNoURLsSpecified
	}
	if len(urls) > maxNumberOfURLsToLoadBalance {
		return nil, ErrTooMuchMultipleURLs
	}

	multipleShortURLs := &MultipleShortURLs{
		Hash:     hashFromURL(strings.Join(urls, "")),
		LongURLs: urls,
	}

	err := b.repository.Save(multipleShortURLs)
	if err != nil {
		return nil, fmt.Errorf("error saving multiple URLs into repository: %w", err)
	}

	return multipleShortURLs, nil
}

func NewLoadBalancer(repository MultipleShortURLsRepository) *LoadBalancer {
	return &LoadBalancer{
		repository: repository,
	}
}
