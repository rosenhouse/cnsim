package distributions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rosenhouse/cnsim/distributions"
)

var _ = Describe("Geometric Distribution with support on the positive integers", func() {
	var (
		dist *distributions.GeometricWithPositiveSupport
	)

	BeforeEach(func() {
		dist = &distributions.GeometricWithPositiveSupport{}
	})

	DescribeTable("sample means",
		func(desiredMean float64) {
			const numSamples = 10000
			var tolerance = 0.05 * desiredMean // prob test suite failure < 0.01
			total := 0
			for i := 0; i < numSamples; i++ {
				sample, err := dist.Sample(desiredMean)
				Expect(err).NotTo(HaveOccurred())
				total += sample
			}
			sampleMean := float64(total) / float64(numSamples)
			Expect(sampleMean).To(BeNumerically("~", desiredMean, tolerance))
		},
		Entry("p=1.0", 1.0),
		Entry("p=0.99", 1.01010101010101),
		Entry("p=0.9", 1.11111111),
		Entry("p=0.8", 1.25),
		Entry("p=0.625", 1.6),
		Entry("p=0.5", 2.0),
		Entry("p=0.4", 2.5),
		Entry("p=0.25", 4.0),
		Entry("p=0.1", 10.0),
		Entry("p=0.01", 100.0),
	)
})
