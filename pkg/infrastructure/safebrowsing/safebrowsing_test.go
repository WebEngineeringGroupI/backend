package safebrowsing_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/safebrowsing"
)

var _ = Describe("SafeBrowsing Validator", func() {
	var (
		validator *safebrowsing.Validator
	)
	BeforeEach(func() {
		var err error
		validator, err = safebrowsing.NewValidator(os.Getenv("SAFE_BROWSING_API_KEY"))
		Expect(err).To(Succeed())
	})

	Context("when checking a safe website", func() {
		It("returns that the URL is valid", func() {
			isSafe, err := validator.ValidateURL("google.com")

			Expect(err).To(Succeed())
			Expect(isSafe).To(BeTrue())
		})
	})

	Context("when checking an unsafe website", func() {
		It("returns that the URL is not valid", func() {
			isSafe, err := validator.ValidateURL("wp.readhere.in")

			Expect(err).To(Succeed())
			Expect(isSafe).To(BeFalse())
		})
	})

	Context("when checking if multiple safe websites are valid", func() {
		It("returns that the URLs are valid", func() {
			isSafe, err := validator.ValidateURLs([]string{"google.com", "youtube.com"})

			Expect(err).To(Succeed())
			Expect(isSafe).To(BeTrue())
		})
	})

	Context("when checking if multiple urls are safe, but one of them is invalid", func() {
		It("returns that the URLs are not valid", func() {
			isSafe, err := validator.ValidateURLs([]string{"google.com", "wp.readhere.in"})

			Expect(err).To(Succeed())
			Expect(isSafe).To(BeFalse())
		})
	})
})
