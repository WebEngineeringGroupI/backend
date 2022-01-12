package safebrowsing_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/safebrowsing"
)

var _ = Describe("SafeBrowsing Validator", func() {
	var (
		validator *safebrowsing.Validator
		ctx       context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
		var err error
		validator, err = safebrowsing.NewValidator(os.Getenv("SAFE_BROWSING_API_KEY"))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when checking if multiple safe websites are valid", func() {
		It("returns that the URLs are valid", func() {
			isSafe, err := validator.ValidateURLs(ctx, []string{"google.com", "youtube.com"})

			Expect(err).ToNot(HaveOccurred())
			Expect(isSafe).To(BeTrue())
		})
	})

	Context("when checking if multiple urls are safe, but one of them is invalid", func() {
		It("returns that the URLs are not valid", func() {
			isSafe, err := validator.ValidateURLs(ctx, []string{"google.com", "wp.readhere.in"})

			Expect(err).ToNot(HaveOccurred())
			Expect(isSafe).To(BeFalse())
		})
	})
})
