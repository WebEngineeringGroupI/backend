package url

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
)

type SingleURLShortener struct {
	repository ShortURLRepository
	metrics    Metrics
	emitter    event.Emitter
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

func (s *SingleURLShortener) HashFromURL(ctx context.Context, aLongURL string) (*ShortURL, error) {
	s.metrics.RecordSingleURLMetrics()

	shortURL := &ShortURL{
		Hash: hashFromURL(aLongURL),
		OriginalURL: OriginalURL{
			URL:     aLongURL,
			IsValid: false,
		},
	}

	err := s.repository.SaveShortURL(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("unable to save shortURL in the repository: %w", err)
	}

	err = s.emitter.EmitShortURLCreated(ctx, shortURL.Hash, shortURL.OriginalURL.URL, shortURL.OriginalURL.IsValid)
	if err != nil {
		return nil, fmt.Errorf("unable to save domain event: %w", err)
	}

	return shortURL, nil
}

func hashFromURL(aLongURL string) string {
	bytes := sha1.Sum([]byte(aLongURL))
	sum := base64.StdEncoding.EncodeToString(bytes[:])
	hash := sum[0:8]
	return hash
}

func NewSingleURLShortener(repository ShortURLRepository, metrics Metrics, emitter event.Emitter) *SingleURLShortener {
	return &SingleURLShortener{
		repository: repository,
		metrics:    metrics,
		emitter:    emitter,
	}
}
