package url_test

import (
	"context"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type FakeFormatter struct {
	longURLs []string
	error    error
}

func (f *FakeFormatter) shouldReturnURLs(longURLs []string) {
	f.longURLs = longURLs
}

func (f *FakeFormatter) shouldReturnError(err error) {
	f.error = err
}

func (f *FakeFormatter) FormatDataToURLs(data []byte) ([]string, error) {
	return f.longURLs, f.error
}

type FakeMetrics struct {
	singleURLMetrics int
	fileURLMetrics   int
	urlsProcessed    int
}

func (f *FakeMetrics) RecordSingleURLMetrics() {
	f.singleURLMetrics++
}

func (f *FakeMetrics) RecordFileURLMetrics() {
	f.fileURLMetrics++
}

func (f *FakeMetrics) RecordUrlsProcessed() {
	f.urlsProcessed++
}

type FakeLoadBalancedURLsRepository struct {
	urls          []*url.LoadBalancedURL
	errorToReturn error
}

func (f *FakeLoadBalancedURLsRepository) shouldReturnError(err error) {
	f.errorToReturn = err
}

func (f *FakeLoadBalancedURLsRepository) SaveLoadBalancedURL(ctx context.Context, urls *url.LoadBalancedURL) error {
	if f.errorToReturn != nil {
		return f.errorToReturn
	}

	f.urls = append(f.urls, urls)
	return nil
}

func (f *FakeLoadBalancedURLsRepository) FindLoadBalancedURLByHash(ctx context.Context, hash string) (*url.LoadBalancedURL, error) {
	panic("implement me")
}
