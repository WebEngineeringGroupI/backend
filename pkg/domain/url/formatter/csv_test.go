package formatter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/formatter"
)

var _ = Describe("CSV Formatter", func() {
	var (
		csvFormatter *formatter.CSV
	)
	BeforeEach(func() {
		csvFormatter = formatter.NewCSV()
	})

	It("transforms a CSV of long URLs into a slice of long URLs", func() {
		longURLs, err := csvFormatter.FormatDataToURLs([]byte("\"https://google.com\"\n\"https://unizar.es\""))

		Expect(err).ToNot(HaveOccurred())
		Expect(longURLs).To(Equal([]string{"https://google.com", "https://unizar.es"}))
	})

	Context("when the CSV is empty", func() {
		It("returns an error", func() {
			longURLs, err := csvFormatter.FormatDataToURLs([]byte(""))

			Expect(err).To(MatchError(url.ErrUnableToConvertDataToLongURLs))
			Expect(longURLs).To(BeNil())
		})
	})
})
