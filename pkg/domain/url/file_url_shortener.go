package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
)

var ErrUnableToConvertDataToLongURLs = errors.New("unable to convert data to long urls")

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type Formatter interface {
	FormatDataToURLs(data []byte) ([]string, error)
}

type FileURLShortener struct {
	repository event.Repository
	formatter  Formatter
	metrics    Metrics
	clock      event.Clock
}

func (s *FileURLShortener) HashesFromURLData(ctx context.Context, data []byte) ([]ShortURL, error) {
	var shortURLs []ShortURL
	s.metrics.RecordFileURLMetrics()

	longURLs, err := s.formatter.FormatDataToURLs(data)
	if err != nil {
		return nil, err
	}

	shortURLs = make([]ShortURL, 0, len(longURLs))
	for _, longURL := range longURLs {
		events := []event.Event{
			&ShortURLCreated{
				Base: event.Base{
					ID:      hashFromURL(longURL),
					Version: 0,
					At:      s.clock.Now(),
				},
				OriginalURL: longURL,
			},
		}
		shortURLs = append(shortURLs, *shortURLFromEvents(events...))
		err = s.repository.Save(ctx, events...)
		if err != nil {
			return nil, fmt.Errorf("unable to save events to repository: %w", err)
		}
	}

	return shortURLs, nil
}

func NewFileURLShortener(repository event.Repository, metrics Metrics, clock event.Clock, formatter Formatter) *FileURLShortener {
	return &FileURLShortener{
		repository: repository,
		formatter:  formatter,
		metrics:    metrics,
		clock:      clock,
	}
}
