package url

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/WebEngineeringGroupI/backend/pkg/domain"
)

type SingleURLShortener struct {
	repository ShortURLRepository
	metrics    Metrics
	outbox     domain.EventOutbox
	clock      Clock
	uuid       UUID
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
func (s *SingleURLShortener) HashFromURL(ctx context.Context, aLongURL string) (*ShortURL, error) {
	s.metrics.RecordSingleURLMetrics()

	shortURL := &ShortURL{
		Hash: hashFromURL(aLongURL),
		OriginalURL: OriginalURL{
			URL:     aLongURL,
			IsValid: false,
		},
	}

	err := s.repository.SaveShortURL(ctx, shortURL) // FIXME(fede): test this error
	if err != nil {
		return nil, fmt.Errorf("unable to save shortURL in the repository: %w", err)
	}

	err = s.outbox.SaveEvent(ctx, NewShortURLCreated(shortURL))
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

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type Clock interface {
	Now() time.Time
}

type UUID interface {
	New() string
}

func NewSingleURLShortener(repository ShortURLRepository, metrics Metrics, outbox domain.EventOutbox, clock Clock, uuid UUID) *SingleURLShortener {
	return &SingleURLShortener{
		repository: repository,
		metrics:    metrics,
		outbox:     outbox,
		clock:      clock,
		uuid:       uuid,
	}
}
