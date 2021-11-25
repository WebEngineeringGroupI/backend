package url_test

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
	singleURLMetrics   int
	multipleURLMetrics int
	fileURLMetrics     int
	urlsProcessed      int
}

func (f *FakeMetrics) RecordSingleURLMetrics() {
	f.singleURLMetrics++
}

func (f *FakeMetrics) RecordMultipleURLMetrics() {
	f.multipleURLMetrics++
}

func (f *FakeMetrics) RecordFileURLMetrics() {
	f.fileURLMetrics++
}

func (f *FakeMetrics) RecordUrlsProcessed() {
	f.urlsProcessed++
}
