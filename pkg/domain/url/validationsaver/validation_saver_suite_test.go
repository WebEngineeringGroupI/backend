package validationsaver_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValidationSaver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ValidationSaver Suite")
}
