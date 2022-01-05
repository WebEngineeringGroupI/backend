package url

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
)

var (
	ErrShortURLNotFound = errors.New("short url not found")
)

type SingleURLShortener struct {
	repository event.Repository
	metrics    Metrics
	clock      event.Clock
}

type OriginalURL struct {
	URL     string
	IsValid bool
}

type ShortURL struct {
	Hash        string
	OriginalURL OriginalURL
	Clicks      int
}

func shortURLFromEvents(events ...event.Event) *ShortURL {
	url := &ShortURL{}
	for _, e := range events {
		_ = url.On(e)
	}
	return url
}

func (s *ShortURL) On(evt event.Event) error {
	switch e := evt.(type) {
	case *ShortURLCreated:
		s.Hash = e.EntityID()
		s.OriginalURL = OriginalURL{URL: e.OriginalURL, IsValid: false}
		s.Clicks = 0
	case *ShortURLVerified:
		s.OriginalURL = OriginalURL{
			URL:     s.OriginalURL.URL,
			IsValid: true,
		}
	case *ShortURLClicked:
		s.Clicks++
	default:
		return event.ErrUnhandledEvent
	}
	return nil
}

var (
	ErrInvalidLongURLSpecified = errors.New("invalid long URL specified")
)

func (s *SingleURLShortener) HashFromURL(ctx context.Context, aLongURL string) (*ShortURL, error) {
	s.metrics.RecordSingleURLMetrics()

	events := []event.Event{
		&ShortURLCreated{
			Base: event.Base{
				ID:      hashFromURL(aLongURL),
				Version: 0,
				At:      s.clock.Now(),
			},
			OriginalURL: aLongURL,
		},
	}

	err := s.repository.Save(ctx, events...)
	if err != nil {
		return nil, fmt.Errorf("unable to save shortURL in the repository: %w", err)
	}

	return shortURLFromEvents(events...), nil
}

func hashFromURL(aLongURL string) string {
	bytes := sha1.Sum([]byte(aLongURL))
	sum := base64.StdEncoding.EncodeToString(bytes[:])
	hash := sum[0:8]
	return hash
}

func NewSingleURLShortener(repository event.Repository, clock event.Clock, metrics Metrics) *SingleURLShortener {
	return &SingleURLShortener{
		repository: repository,
		clock:      clock,
		metrics:    metrics,
	}
}
