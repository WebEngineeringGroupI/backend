package event

import (
	"time"
)

type ShortURLCreated struct {
	EventID     string
	Creation    time.Time
	Hash        string
	OriginalURL string
}

func (s *ShortURLCreated) ID() string {
	return s.EventID
}

func (s *ShortURLCreated) HappenedOn() time.Time {
	return s.Creation
}

func NewShortURLCreated(id string, happenedOn time.Time, hash string, originalURL string) *ShortURLCreated {
	return &ShortURLCreated{
		EventID:     id,
		Creation:    happenedOn,
		Hash:        hash,
		OriginalURL: originalURL,
	}
}
