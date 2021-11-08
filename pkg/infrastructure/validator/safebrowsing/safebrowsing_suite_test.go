package safebrowsing_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSafebrowsing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Safebrowsing Suite")
}
