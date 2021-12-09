package url

import (
	"context"
	"errors"
)

var (
	ErrShortURLNotFound = errors.New("short url not found")
)

type ShortURLRepository interface {
	FindShortURLByHash(ctx context.Context, hash string) (*ShortURL, error)
	SaveShortURL(ctx context.Context, url *ShortURL) error
}
