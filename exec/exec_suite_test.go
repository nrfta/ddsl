package exec

import (
	"testing"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExec(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Exec Suite")
}
