package redirect

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type LoadBalancerRedirectorService struct {
	repository event.Repository
}

func (r *LoadBalancerRedirectorService) ReturnAValidOriginalURL(ctx context.Context, hash string) (string, error) {
	loadBalancedURLsEntity, _, err := r.repository.Load(ctx, hash)
	if errors.Is(err, event.ErrEntityNotFound) {
		return "", url.ErrValidURLNotFound
	}

	if err != nil {
		return "", err
	}
	loadBalancedURLs, ok := loadBalancedURLsEntity.(*url.LoadBalancedURL)
	if !ok {
		return "", fmt.Errorf("the entity loaded is not a LoadBalancedURL: %w", url.ErrValidURLNotFound)
	}

	validURLs := r.filterValidURLs(loadBalancedURLs.LongURLs)
	if len(validURLs) == 0 {
		return "", url.ErrValidURLNotFound
	}

	nextIndex := rand.Intn(len(validURLs)) // TODO(fede): Implement roundRobin with MemCached service
	return validURLs[nextIndex], nil
}

func (r *LoadBalancerRedirectorService) filterValidURLs(originalURLs []url.OriginalURL) []string {
	validURLs := []string{}
	for _, aURL := range originalURLs {
		if aURL.IsValid {
			validURLs = append(validURLs, aURL.URL)
		}
	}
	return validURLs
}

func NewLoadBalancerRedirectorService(repository event.Repository) *LoadBalancerRedirectorService {
	return &LoadBalancerRedirectorService{repository: repository}
}
