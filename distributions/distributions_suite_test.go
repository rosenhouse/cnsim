package distributions_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDistributions(t *testing.T) {
	rand.Seed(config.GinkgoConfig.RandomSeed + int64(GinkgoParallelNode()))

	RegisterFailHandler(Fail)
	RunSpecs(t, "Distributions Suite")
}
