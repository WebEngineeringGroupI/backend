package url_test

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type FakeURLValidator struct {
	returnValidURL bool
	returnError    error
}

func (f *FakeURLValidator) shouldReturnValidURL(validURL bool) {
	f.returnValidURL = validURL
}

func (f *FakeURLValidator) shouldReturnError(err error) {
	f.returnError = err
}

func (f *FakeURLValidator) ValidateURL(url string) (bool, error) {
	return f.returnValidURL, f.returnError
}

func (f *FakeURLValidator) ValidateURLs(url []string) (bool, error) {
	return f.returnValidURL, f.returnError
}

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

type FakeMultipleShortURLsRepository struct {
	urls          []*url.LoadBalancedURL
	errorToReturn error
}

func (f *FakeMultipleShortURLsRepository) shouldReturnError(err error) {
	f.errorToReturn = err
}

func (f *FakeMultipleShortURLsRepository) Save(urls *url.LoadBalancedURL) error {
	if f.errorToReturn != nil {
		return f.errorToReturn
	}

	f.urls = append(f.urls, urls)
	return nil
}

func (f *FakeMultipleShortURLsRepository) FindByHash(hash string) (*url.LoadBalancedURL, error) {
	panic("implement me")
}
