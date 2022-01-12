package redirector_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRedirector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redirector Suite")
}
