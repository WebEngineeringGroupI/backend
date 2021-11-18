package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain"
)

var _ = Describe("Domain / Whole Url", func() {
	When("it's provided a hash", func() {
		It("returns the whole url", func() {
			wholeURL := domain.NewWholeURL("http://foobar.baz")
			url := wholeURL.FromHash("someHashHere")

			Expect(url).To(Equal("http://foobar.baz/r/someHashHere"))
		})
		It("returns the whole url even if created with trailing slash", func() {
			wholeURL := domain.NewWholeURL("http://foobar.baz/")
			url := wholeURL.FromHash("someHashHere")

			Expect(url).To(Equal("http://foobar.baz/r/someHashHere"))
		})
	})
})
