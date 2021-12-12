package inmemory

import (
	"context"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Repository struct {
	shortURLs        map[string]*url.ShortURL
	loadBalancedURLs map[string]*url.LoadBalancedURL
	eventOutbox      map[string]event.Event
}

func (f *Repository) SaveEvent(ctx context.Context, event event.Event) error {
	f.eventOutbox[event.ID()] = event
	return nil
}

func (f *Repository) SaveShortURL(ctx context.Context, url *url.ShortURL) error {
	f.shortURLs[url.Hash] = url
	return nil
}

func (f *Repository) FindShortURLByHash(ctx context.Context, hash string) (*url.ShortURL, error) {
	shortURL, ok := f.shortURLs[hash]
	if !ok {
		return nil, url.ErrShortURLNotFound
	}

	return shortURL, nil
}

func (f *Repository) SaveLoadBalancedURL(ctx context.Context, urls *url.LoadBalancedURL) error {
	f.loadBalancedURLs[urls.Hash] = urls
	return nil
}

func (f *Repository) FindLoadBalancedURLByHash(ctx context.Context, hash string) (*url.LoadBalancedURL, error) {
	loadBalancedURL, ok := f.loadBalancedURLs[hash]
	if !ok {
		return nil, url.ErrValidURLNotFound // FIXME(fede): We should return other kind of error?
	}
	return loadBalancedURL, nil
}

func NewRepository() *Repository {
	return &Repository{
		shortURLs:        map[string]*url.ShortURL{},
		loadBalancedURLs: map[string]*url.LoadBalancedURL{},
		eventOutbox:      map[string]event.Event{},
	}
}
