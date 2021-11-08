package pipeline_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMultiple(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multiple Suite")
}
