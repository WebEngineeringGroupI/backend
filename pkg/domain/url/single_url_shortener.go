package url

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
)

type SingleURLShortener struct {
	repository ShortURLRepository
	validator  Validator
	metrics    Metrics
}

type OriginalURL struct {
	URL     string
	IsValid bool
}

type ShortURL struct {
	Hash        string
	OriginalURL OriginalURL
}

var (
	ErrInvalidLongURLSpecified = errors.New("invalid long URL specified")
)

// FIXME(fede): Rename to something like ShortURLFromLong
func (s *SingleURLShortener) HashFromURL(aLongURL string) (*ShortURL, error) {
	s.metrics.RecordSingleURLMetrics()
	isValidURL, err := s.validator.ValidateURLs([]string{aLongURL})
	if err != nil {
		return nil, err
	}
	if !isValidURL {
		return nil, ErrInvalidLongURLSpecified
	}

	shortURL := &ShortURL{
		Hash: hashFromURL(aLongURL),
		OriginalURL: OriginalURL{
			URL:     aLongURL,
			IsValid: true,
		},
	}

	err = s.repository.SaveShortURL(shortURL) // FIXME(fede): test this error
	if err != nil {
		return nil, fmt.Errorf("unable to save shortURL in the repository: %w", err)
	}

	return shortURL, nil
}

func hashFromURL(aLongURL string) string {
	bytes := sha1.Sum([]byte(aLongURL))
	sum := base64.StdEncoding.EncodeToString(bytes[:])
	hash := sum[0:8]
	return hash
}

func NewSingleURLShortener(repository ShortURLRepository, validator Validator, metrics Metrics) *SingleURLShortener {
	return &SingleURLShortener{
		repository: repository,
		validator:  validator,
		metrics:    metrics,
	}
}
