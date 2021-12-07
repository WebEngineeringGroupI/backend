package redirect

import (
	"errors"
	"fmt"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"math/rand"
)

var (
	ValidURLNotFound = errors.New("valid URL not found")
)

type LoadBalancedURLsRepository interface {
	FindByHash(hash string) (*url.LoadBalancedURL, error)
}

type LoadBalancerRedirector struct {
	repository LoadBalancedURLsRepository
}

func (r *LoadBalancerRedirector) ReturnAValidOriginalURL(hash string) (string, error) {
	loadBalancedURLs, err := r.repository.FindByHash(hash)
	if err != nil {
		return "", err
	}

	validURLs := r.filterValidURLs(loadBalancedURLs.LongURLs)
	if len(validURLs) == 0 {
		return "", fmt.Errorf("there are no valid URLs to redirect to")
	}

	randomIndex := rand.Intn(len(validURLs))
	return validURLs[randomIndex], nil
}

func (r *LoadBalancerRedirector) filterValidURLs(originalURLs []url.OriginalURL) []string {
	validURLs := []string{}
	for _, aURL := range originalURLs {
		if aURL.IsValid {
			validURLs = append(validURLs, aURL.URL)
		}
	}
	return validURLs
}

func NewLoadBalancerRedirector(repository LoadBalancedURLsRepository) *LoadBalancerRedirector {
	return &LoadBalancerRedirector{repository: repository}
}
