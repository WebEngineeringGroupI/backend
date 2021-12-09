package url

import (
	"context"
	"errors"
	"fmt"
)

var ErrUnableToConvertDataToLongURLs = errors.New("unable to convert data to long urls")

//go:generate mockgen -source=$GOFILE -destination=./mocks/${GOFILE} -package=mocks
type Formatter interface {
	FormatDataToURLs(data []byte) ([]string, error)
}

type FileURLShortener struct {
	repository ShortURLRepository
	formatter  Formatter
	metrics    Metrics
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
		shortURL := ShortURL{
			Hash: hashFromURL(longURL),
			OriginalURL: OriginalURL{
				URL:     longURL,
				IsValid: false,
			},
		}

		err := s.repository.SaveShortURL(ctx, &shortURL)
		if err != nil {
			return nil, fmt.Errorf("unable to save URL '%s' to repository: %w", longURL, err)
		}

		shortURLs = append(shortURLs, shortURL)
	}

	return shortURLs, nil
}

func NewFileURLShortener(repository ShortURLRepository, metrics Metrics, formatter Formatter) *FileURLShortener {
	return &FileURLShortener{
		repository: repository,
		formatter:  formatter,
		metrics:    metrics,
	}
}
