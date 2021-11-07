package url

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
)

type Shortener struct {
	repository ShortURLRepository
	validator  Validator
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
	isValidURL, err := s.validator.ValidateURL(aLongURL)
	if err != nil {
		return nil, err
	}
	if !isValidURL {
		return nil, ErrInvalidLongURLSpecified
	}

	bytes := sha1.Sum([]byte(aLongURL))
	sum := base64.StdEncoding.EncodeToString(bytes[:])

	shortURL := &ShortURL{
		Hash:    sum[0:8],
		LongURL: aLongURL,
	}

	err = s.repository.Save(shortURL) // FIXME(fede): test this error
	if err != nil {
		return nil, fmt.Errorf("unable to save shortURL in the repository: %w", err)
	}

	return shortURL, nil
}

func NewShortener(repository ShortURLRepository, validator Validator) *Shortener {
	return &Shortener{
		repository: repository,
		validator:  validator,
	}
}
