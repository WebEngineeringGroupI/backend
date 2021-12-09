package url

import (
	"time"
)

type ShortURLCreated struct {
	id         string
	happenedOn time.Time
	url        *ShortURL
}

func (s *ShortURLCreated) ID() string {
	return s.id
}

func (s *ShortURLCreated) HappenedOn() time.Time {
	return s.happenedOn
}

func (s *ShortURLCreated) ShortURL() *ShortURL {
	return s.url
}

func NewShortURLCreated(aURL *ShortURL) *ShortURLCreated {
	return &ShortURLCreated{
		url: aURL,
	}
}
