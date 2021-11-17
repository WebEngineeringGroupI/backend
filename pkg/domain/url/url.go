package url

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type Shortener struct {
	repository ShortURLRepository
	metrics Metrics
}

type ShortURL struct {
	Hash    string
	LongURL string
}

var (
	ErrInvalidLongURLSpecified = errors.New("invalid long URL specified")
)

// FIXME(fede): Rename to something like ShortURLFromLong
func (s *Shortener) HashFromURL(aLongURL string) (*ShortURL, error) {
	if !strings.HasPrefix(aLongURL, "http://") && !strings.HasPrefix(aLongURL, "https://") {
		return nil, ErrInvalidLongURLSpecified
	}

	bytes := sha1.Sum([]byte(aLongURL))
	sum := base64.StdEncoding.EncodeToString(bytes[:])

	shortURL := &ShortURL{
		Hash:    sum[0:8],
		LongURL: aLongURL,
	}

	err := s.repository.Save(shortURL) // FIXME(fede): test this error
	if err != nil {
		return nil, fmt.Errorf("unable to save shortURL in the repository: %w", err)
	}

	s.metrics.RecordUrlShorted()
	return shortURL, nil
}

func NewShortener(repository ShortURLRepository) *Shortener {
	return &Shortener{
		repository: repository,
	}
}
