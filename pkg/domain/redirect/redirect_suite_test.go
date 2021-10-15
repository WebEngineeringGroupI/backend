package redirect_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRedirect(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redirect Suite")
}
