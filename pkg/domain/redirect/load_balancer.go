package redirect

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

var (
	ErrValidURLNotFound = errors.New("valid URL not found")
)

type LoadBalancerRedirector struct {
	repository url.LoadBalancedURLsRepository
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

func NewLoadBalancerRedirector(repository url.LoadBalancedURLsRepository) *LoadBalancerRedirector {
	return &LoadBalancerRedirector{repository: repository}
}
