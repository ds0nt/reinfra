package reinfra_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestReinfra(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reinfra Suite")
}
