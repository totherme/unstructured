package unstructured_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUnstructured(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unstructured Suite")
}
