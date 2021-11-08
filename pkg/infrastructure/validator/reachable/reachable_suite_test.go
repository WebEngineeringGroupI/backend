package reachable_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestReachable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reachable Suite")
}
