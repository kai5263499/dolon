package dolon_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDolon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dolon Suite")
}
