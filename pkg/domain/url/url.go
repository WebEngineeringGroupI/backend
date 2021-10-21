package url

import (
	`crypto/sha1`
	`encoding/base64`
	`strings`
)

type Shortener struct {
	repository ShortURLRepository
}

type ShortURL struct {
	Hash    string
	LongURL string
}

// FIXME(fede): Rename to something like ShortURLFromLong
func (s *Shortener) HashFromURL(aLongURL string) *ShortURL {
	if !strings.HasPrefix(aLongURL, "http://") && !strings.HasPrefix(aLongURL, "https://") {
		return nil
	}

	bytes := sha1.Sum([]byte(aLongURL))
	sum := base64.StdEncoding.EncodeToString(bytes[:])

	shortURL := &ShortURL{
		Hash:    sum[0:8],
		LongURL: aLongURL,
	}

	s.repository.Save(shortURL)
	return shortURL
}

func NewShortener(repository ShortURLRepository) *Shortener {
	return &Shortener{
		repository: repository,
	}
}
