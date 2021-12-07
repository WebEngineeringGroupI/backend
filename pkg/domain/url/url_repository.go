package url

import (
	"errors"
)

var (
	ErrShortURLNotFound = errors.New("short url not found")
)

type ShortURLRepository interface {
	FindShortURLByHash(hash string) (*ShortURL, error)
	SaveShortURL(url *ShortURL) error
}
