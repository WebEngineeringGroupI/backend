package url_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Multiple URL Shortener", func() {
	var (
		shortener  *url.FileURLShortener
		repository url.ShortURLRepository
		formatter  *FakeFormatter
		validator  *FakeURLValidator
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		validator = &FakeURLValidator{returnValidURL: true}
		formatter = &FakeFormatter{}
		formatter.shouldReturnURLs(aLongURLSet())
		shortener = url.NewFileURLShortener(repository, validator, formatter)
	})

	Context("when providing multiple long URLs", func() {
		It("generates a hash for each one", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())

			Expect(err).To(Succeed())
			Expect(shortURLs).To(HaveLen(2))
			Expect(shortURLs[0].Hash).To(HaveLen(8))
			Expect(shortURLs[1].Hash).To(HaveLen(8))
		})

		It("contains the real values from the original URLs", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())

			Expect(err).To(Succeed())
			Expect(shortURLs[0].LongURL).To(Equal("https://google.com"))
			Expect(shortURLs[1].LongURL).To(Equal("https://unizar.es"))
		})

		It("generates different short URL hashes for each of the long URLs", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())

			Expect(err).To(Succeed())
			Expect(shortURLs[0].Hash).ToNot(Equal(shortURLs[1].Hash))
		})

		Context("but the provided data is not valid", func() {
			It("returns the error since it's unable to transform the data", func() {
				formatter.shouldReturnError(errors.New("unknown error"))
				shortURLs, err := shortener.HashesFromURLData(aLongURLData())

				Expect(err).To(MatchError("unknown error"))
				Expect(shortURLs).To(BeNil())
			})
		})

		Context("and the provided URL is not valid", func() {
			It("validates that the provided URL is not valid", func() {
				validator.shouldReturnValidURL(false)
				shortURLs, err := shortener.HashesFromURLData(aLongURLData())

				Expect(err).To(MatchError(url.ErrInvalidLongURLSpecified))
				Expect(shortURLs).To(BeNil())
			})
		})

		Context("but the validator returns an error", func() {
			It("returns the error since it's unable to validate the URL", func() {
				validator.shouldReturnError(errors.New("unknown error"))
				shortURLs, err := shortener.HashesFromURLData(aLongURLData())

				Expect(err).To(MatchError("unknown error"))
				Expect(shortURLs).To(BeNil())
			})
		})

		It("stores the short URL in a repository", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())
			Expect(err).To(Succeed())

			firstURL, err := repository.FindByHash(shortURLs[0].Hash)
			Expect(err).To(Succeed())
			secondURL, err := repository.FindByHash(shortURLs[1].Hash)
			Expect(err).To(Succeed())

			Expect(firstURL.Hash).To(Equal(shortURLs[0].Hash))
			Expect(firstURL.LongURL).To(Equal("https://google.com"))
			Expect(secondURL.Hash).To(Equal(shortURLs[1].Hash))
			Expect(secondURL.LongURL).To(Equal("https://unizar.es"))
		})
	})
})

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

func aLongURLData() []byte {
	return []byte(`"https://google.com"
"https://unizar.es"`)
}

func aLongURLSet() []string {
	return []string{"https://google.com", "https://unizar.es"}
}
