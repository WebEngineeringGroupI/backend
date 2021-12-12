package event

import (
	"time"
)

type LoadBalancedURLCreated struct {
	EventID      string
	Creation     time.Time
	Hash         string
	OriginalURLs []string
}

func (s *LoadBalancedURLCreated) ID() string {
	return s.EventID
}

func (s *LoadBalancedURLCreated) HappenedOn() time.Time {
	return s.Creation
}

func NewLoadBalancedURLCreated(id string, happenedOn time.Time, hash string, originalURLs []string) *LoadBalancedURLCreated {
	return &LoadBalancedURLCreated{
		EventID:      id,
		Creation:     happenedOn,
		Hash:         hash,
		OriginalURLs: originalURLs,
	}
}
