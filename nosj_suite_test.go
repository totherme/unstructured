package nosj_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNosj(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nosj Suite")
}
