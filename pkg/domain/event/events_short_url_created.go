package event

import (
	"time"
)

type ShortURLCreated struct {
	EventID     string
	Creation    time.Time
	Hash        string
	OriginalURL string
	IsValid     bool
}

func (s *ShortURLCreated) ID() string {
	return s.EventID
}

func (s *ShortURLCreated) HappenedOn() time.Time {
	return s.Creation
}

func NewShortURLCreated(id string, happenedOn time.Time, hash string, originalURL string, isValid bool) *ShortURLCreated {
	return &ShortURLCreated{
		EventID:     id,
		Creation:    happenedOn,
		Hash:        hash,
		OriginalURL: originalURL,
		IsValid:     isValid,
	}
}
