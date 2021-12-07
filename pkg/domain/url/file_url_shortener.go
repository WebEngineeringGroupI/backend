package url

import (
	"errors"
	"fmt"
)

var ErrUnableToConvertDataToLongURLs = errors.New("unable to convert data to long urls")

type Formatter interface {
	FormatDataToURLs(data []byte) ([]string, error)
}

type FileURLShortener struct {
	repository ShortURLRepository
	validator  Validator
	formatter  Formatter
	metrics    Metrics
}

func (s *FileURLShortener) HashesFromURLData(data []byte) ([]ShortURL, error) {
	var shortURLs []ShortURL
	s.metrics.RecordFileURLMetrics()

	longURLs, err := s.formatter.FormatDataToURLs(data)
	if err != nil {
		return nil, err
	}

	urlsAreValid, err := s.validator.ValidateURLs(longURLs)
	if err != nil {
		return nil, err
	}
	if !urlsAreValid {
		return nil, ErrInvalidLongURLSpecified
	}

	shortURLs = make([]ShortURL, 0, len(longURLs))
	for _, longURL := range longURLs {
		s.metrics.RecordUrlsProcessed()
		shortURL := ShortURL{
			Hash: hashFromURL(longURL),
			OriginalURL: OriginalURL{
				URL:     longURL,
				IsValid: true,
			},
		}

		err := s.repository.Save(&shortURL)
		if err != nil {
			return nil, fmt.Errorf("unable to save URL '%s' to repository: %w", longURL, err)
		}

		shortURLs = append(shortURLs, shortURL)
	}

	return shortURLs, nil
}

func NewFileURLShortener(repository ShortURLRepository, validator Validator, metrics Metrics, formatter Formatter) *FileURLShortener {
	return &FileURLShortener{
		repository: repository,
		validator:  validator,
		formatter:  formatter,
		metrics:    metrics,
	}
}
