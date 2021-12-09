package schema_test

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/schema"
)

var _ = Describe("Schema Validator", func() {
	var (
		ctx context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
	})
	It("validates a URL that contains the http schema", func() {
		validator := schema.NewValidator("http", "https")
		isValid, err := validator.ValidateURLs(ctx, []string{"http://google.com"})

		Expect(err).ToNot(HaveOccurred())
		Expect(isValid).To(BeTrue())
	})
	It("validates a URL that contains the https schema", func() {
		validator := schema.NewValidator("http", "https")
		isValid, err := validator.ValidateURLs(ctx, []string{"https://google.com"})

		Expect(err).ToNot(HaveOccurred())
		Expect(isValid).To(BeTrue())
	})
	It("fails to validate a URL that contains the ftp schema", func() {
		validator := schema.NewValidator("http", "https")
		isValid, err := validator.ValidateURLs(ctx, []string{"ftp://google.com"})

		Expect(err).ToNot(HaveOccurred())
		Expect(isValid).To(BeFalse())
	})
	It("validates a random schema", func() {
		randomSchema := randomStringOfLength(5)
		validator := schema.NewValidator(randomSchema)
		isValid, err := validator.ValidateURLs(ctx, []string{fmt.Sprintf("%s://google.com", randomSchema)})

		Expect(err).ToNot(HaveOccurred())
		Expect(isValid).To(BeTrue())
	})

})

func randomStringOfLength(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
