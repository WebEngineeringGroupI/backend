package click_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestClick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Click Suite")
}
