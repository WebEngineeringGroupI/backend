package url

import (
	"errors"
)

var (
	ErrShortURLNotFound = errors.New("short url not found")
)

type ShortURLRepository interface {
	FindByHash(hash string) (*ShortURL, error)
	Save(url *ShortURL) error
}

type Metrics interface {
	RecordUrlShorted()
}
