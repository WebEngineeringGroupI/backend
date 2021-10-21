package click_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Click Suite")
}
