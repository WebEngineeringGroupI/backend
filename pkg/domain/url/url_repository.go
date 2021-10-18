package url

import (
	`errors`
)

var (
	ErrShortURLNotFound = errors.New("short url not found")
)

type ShortURLRepository interface {
	FindByHash(hash string) *ShortURL
	Save(url *ShortURL)
}
