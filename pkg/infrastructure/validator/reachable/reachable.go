package reachable

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Validator struct {
	client     *http.Client
	maxTimeout time.Duration
}

type invalidURL struct {
	url string
	err error
}

func (v *Validator) ValidateURLs(ctx context.Context, urls []string) (bool, error) {
	validURLCh := make(chan string, len(urls))
	invalidURLCh := make(chan invalidURL, len(urls))

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))
	for _, url := range urls {
		go v.validateURL(ctx, url, wg, validURLCh, invalidURLCh)
	}
	wg.Wait()

	for i := 0; i < len(urls); i++ {
		select {
		case <-validURLCh:
		case invalid := <-invalidURLCh:
			return false, fmt.Errorf("%w: %s", url.ErrUnableToValidateURLs, invalid.err)
		}
	}
	return true, nil
}

func (v *Validator) validateURL(ctx context.Context, url string, wg *sync.WaitGroup, validURLCh chan<- string, invalidURLCh chan<- invalidURL) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(ctx, v.maxTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		invalidURLCh <- invalidURL{url, err}
		return
	}

	response, err := v.client.Do(request)
	if err != nil {
		invalidURLCh <- invalidURL{url, err}
		return
	}

	if response.StatusCode != http.StatusOK {
		invalidURLCh <- invalidURL{url, fmt.Errorf("could not reach URL '%s': '%s'", url, response.Status)}
		return
	}

	validURLCh <- url
}

func NewValidator(client *http.Client, maxTimeout time.Duration) *Validator {
	return &Validator{
		client:     client,
		maxTimeout: maxTimeout,
	}
}
