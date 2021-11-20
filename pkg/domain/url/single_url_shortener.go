package url

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
)

type SingleURLShortener struct {
	repository    ShortURLRepository
	validator     Validator
	recordMetrics RecordMetrics
}

type ShortURL struct {
	Hash    string
	LongURL string
}

var (
	ErrInvalidLongURLSpecified = errors.New("invalid long URL specified")
)

// FIXME(fede): Rename to something like ShortURLFromLong
func (s *SingleURLShortener) HashFromURL(aLongURL string) (*ShortURL, error) {
	isValidURL, err := s.validator.ValidateURLs([]string{aLongURL})
	if err != nil {
		return nil, err
	}
	if !isValidURL {
		return nil, ErrInvalidLongURLSpecified
	}

	shortURL := &ShortURL{
		Hash:    hashFromURL(aLongURL),
		LongURL: aLongURL,
	}

	err = s.repository.Save(shortURL) // FIXME(fede): test this error
	if err != nil {
		return nil, fmt.Errorf("unable to save shortURL in the repository: %w", err)
	}

	// Increment number of Urls processed
	s.recordMetrics.RecordUrlsProcessed()
	// Increment number of SingleUrls processed
	s.recordMetrics.RecordSingleURLMetrics()

	return shortURL, nil
}

func hashFromURL(aLongURL string) string {
	bytes := sha1.Sum([]byte(aLongURL))
	sum := base64.StdEncoding.EncodeToString(bytes[:])
	hash := sum[0:8]
	return hash
}

func NewSingleURLShortener(repository ShortURLRepository, validator Validator) *SingleURLShortener {
	return &SingleURLShortener{
		repository: repository,
		validator:  validator,
	}
}
