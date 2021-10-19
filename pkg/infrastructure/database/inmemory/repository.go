package inmemory

import (
	`github.com/WebEngineeringGroupI/backend/pkg/domain/url`
)

type Repository struct {
	urls map[string]*url.ShortURL
}

func (f *Repository) Save(url *url.ShortURL) error {
	f.urls[url.Hash] = url
	return nil
}

func (f *Repository) FindByHash(hash string) (*url.ShortURL, error) {
	shortURL, ok := f.urls[hash]
	if !ok {
		return nil, url.ErrShortURLNotFound
	}

	return shortURL, nil
}

func NewRepository() *Repository {
	return &Repository{
		urls: map[string]*url.ShortURL{},
	}
}
