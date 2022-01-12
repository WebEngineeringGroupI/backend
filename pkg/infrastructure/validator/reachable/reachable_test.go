package reachable_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/reachable"
)

var _ = Describe("Reachable", func() {
	var (
		validator *reachable.Validator
		ctx       context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
		validator = reachable.NewValidator(http.DefaultClient, 2*time.Second)
	})

	It("validates correctly a reachable URL", func() {
		done := make(chan interface{})
		go func() {
			defer GinkgoRecover()
			isValid, err := validator.ValidateURLs(ctx, validURLsToValidate())

			Expect(err).ToNot(HaveOccurred())
			Expect(isValid).To(BeTrue())
			close(done)
		}()
		Eventually(done, 3.0).Should(BeClosed())
	})

	When("a URL is not reachable", func() {
		It("fails with an error", func() {
			done := make(chan interface{})
			go func() {
				defer GinkgoRecover()
				isValid, err := validator.ValidateURLs(ctx, unreachableURLsToValidate())

				Expect(err).To(MatchError(url.ErrUnableToValidateURLs))
				Expect(err.Error()).To(ContainSubstring("connection refused"))
				Expect(isValid).To(BeFalse())
				close(done)
			}()
			Eventually(done, 3.0).Should(BeClosed())
		})
	})
	When("a URL returns an HTTP code different from 200", func() {
		It("fails with an error", func() {
			done := make(chan interface{})
			go func() {
				defer GinkgoRecover()
				isValid, err := validator.ValidateURLs(ctx, errorURLsToValidate())

				Expect(err).To(MatchError(url.ErrUnableToValidateURLs))
				Expect(err.Error()).To(ContainSubstring("404 Not Found"))
				Expect(isValid).To(BeFalse())
				close(done)
			}()
			Eventually(done, 3.0).Should(BeClosed())
		})
	})
	When("a URL has an invalid certificate", func() {
		It("fails with an error", func() {
			done := make(chan interface{})
			go func() {
				defer GinkgoRecover()
				isValid, err := validator.ValidateURLs(ctx, invalidCertificateURLsToValidate())

				Expect(err).To(MatchError(url.ErrUnableToValidateURLs))
				Expect(err.Error()).To(ContainSubstring("certificate signed by unknown authority"))
				Expect(isValid).To(BeFalse())
				close(done)
			}()
			Eventually(done, 3.0).Should(BeClosed())
		})
	})
})

func validURLsToValidate() []string {
	return []string{
		"https://google.es",
		"https://youtube.com",
		"http://unizar.es",
		"https://www.amazon.com/",
		"https://www.reddit.com",
	}
}

func unreachableURLsToValidate() []string {
	return append(validURLsToValidate(), "http://localhost:21305")
}

func errorURLsToValidate() []string {
	return append(validURLsToValidate(), "http://unizar.es/non-existing-path")
}

func invalidCertificateURLsToValidate() []string {
	return append(validURLsToValidate(), "https://self-signed.badssl.com/")
}
