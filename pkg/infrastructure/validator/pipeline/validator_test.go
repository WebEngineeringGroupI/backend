package pipeline_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/pipeline"
)

var _ = Describe("Multiple Validator", func() {
	var (
		validator    *pipeline.Validator
		validatorOne *FakeValidator
		validatorTwo *FakeValidator
	)
	BeforeEach(func() {
		validatorOne = &FakeValidator{valueToReturn: true}
		validatorTwo = &FakeValidator{valueToReturn: true}
		validator = pipeline.NewValidator(validatorOne, validatorTwo)
	})

	It("executes all validators in the pipeline", func() {
		validURLs, err := validator.ValidateURLs([]string{"google.com"})

		Expect(err).ToNot(HaveOccurred())
		Expect(validURLs).To(BeTrue())
		Expect(validatorOne.shouldHaveBeenCalled()).To(BeTrue())
		Expect(validatorTwo.shouldHaveBeenCalled()).To(BeTrue())
	})

	Context("when the first validator in the pipeline fails to validate", func() {
		It("stops the validation and the second is not executed due to invalid URL", func() {
			validatorOne.shouldReturn(false)
			validURLs, err := validator.ValidateURLs([]string{"google.com"})

			Expect(err).ToNot(HaveOccurred())
			Expect(validURLs).To(BeFalse())
			Expect(validatorOne.shouldHaveBeenCalled()).To(BeTrue())
			Expect(validatorTwo.shouldHaveBeenCalled()).To(BeFalse())
		})
		It("stops the validation and the second is not executed due to error in execution", func() {
			validatorOne.shouldErrorWith(errors.New("unknown error"))
			validURLs, err := validator.ValidateURLs([]string{"google.com"})

			Expect(err).To(MatchError("unknown error"))
			Expect(validURLs).To(BeFalse())
			Expect(validatorOne.shouldHaveBeenCalled()).To(BeTrue())
			Expect(validatorTwo.shouldHaveBeenCalled()).To(BeFalse())
		})
	})

	Context("when the second validator in the pipeline fails to validate", func() {
		It("returns the result of the second validator", func() {
			validatorTwo.shouldReturn(false)
			validURLs, err := validator.ValidateURLs([]string{"google.com"})

			Expect(err).ToNot(HaveOccurred())
			Expect(validURLs).To(BeFalse())
			Expect(validatorOne.shouldHaveBeenCalled()).To(BeTrue())
			Expect(validatorTwo.shouldHaveBeenCalled()).To(BeTrue())
		})
		It("stops the validation and due to error in execution", func() {
			validatorTwo.shouldErrorWith(errors.New("unknown error"))
			validURLs, err := validator.ValidateURLs([]string{"google.com"})

			Expect(err).To(MatchError("unknown error"))
			Expect(validURLs).To(BeFalse())
			Expect(validatorOne.shouldHaveBeenCalled()).To(BeTrue())
			Expect(validatorTwo.shouldHaveBeenCalled()).To(BeTrue())
		})
	})
})

type FakeValidator struct {
	hasBeenCalled bool
	valueToReturn bool
	errorToReturn error
}

func (f *FakeValidator) shouldReturn(value bool) {
	f.valueToReturn = value
}

func (f *FakeValidator) shouldHaveBeenCalled() bool {
	return f.hasBeenCalled
}
func (f *FakeValidator) shouldErrorWith(err error) {
	f.errorToReturn = err
}

func (f *FakeValidator) ValidateURLs(url []string) (bool, error) {
	f.hasBeenCalled = true
	return f.valueToReturn, f.errorToReturn
}
