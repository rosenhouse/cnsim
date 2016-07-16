package simulate_test

import (
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/rosenhouse/cnsim/models"
	"github.com/rosenhouse/cnsim/simulate"
)

var _ = Describe("Steady state simulator", func() {
	var (
		sim    *simulate.SteadyState
		logger *lagertest.TestLogger
		req    models.SteadyStateRequest
	)

	BeforeEach(func() {
		sim = &simulate.SteadyState{}
		logger = lagertest.NewTestLogger("test")
		req = models.SteadyStateRequest{
			NumHosts:            1000,
			NumApps:             10000,
			MeanInstancesPerApp: 50,
		}
	})

	Describe("Execute", func() {
		It("logs on start and stop", func() {
			sim.Execute(logger, req)
			Expect(len(logger.LogMessages())).To(BeNumerically(">=", 2))
		})

		It("logs the structured request and responses", func() {
			sim.Execute(logger, req)

			Expect(logger.Buffer()).To(gbytes.Say(`start.*input.*1000`))
			Expect(logger.Buffer()).To(gbytes.Say(`success.*output.*500`))
		})

		It("returns the request data along with the response", func() {
			resp, err := sim.Execute(logger, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Request).To(Equal(req))
		})

		It("computes the average instances per host", func() {
			resp, err := sim.Execute(logger, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.MeanInstancesPerHost).To(Equal(500.0))
		})
	})

	Describe("Validate", func() {
		It("returns nil when values are within their allowed ranges", func() {
			Expect(sim.Validate(req)).To(Succeed())
		})
		It("returns an error when hosts out of range", func() {
			var bad models.SteadyStateRequest
			bad = req
			bad.NumHosts = 0
			Expect(sim.Validate(bad)).To(MatchError("NumHosts must be 1 - 1000"))

			bad = req
			bad.NumHosts = 1001
			Expect(sim.Validate(bad)).To(MatchError("NumHosts must be 1 - 1000"))

			bad = req
			bad.NumApps = 0
			Expect(sim.Validate(bad)).To(MatchError("NumApps must be 1 - 65534"))

			bad = req
			bad.NumApps = 65535
			Expect(sim.Validate(bad)).To(MatchError("NumApps must be 1 - 65534"))

			bad = req
			bad.MeanInstancesPerApp = 0
			Expect(sim.Validate(bad)).To(MatchError("MeanInstancesPerApp must be 1 - 100"))

			bad = req
			bad.MeanInstancesPerApp = 101
			Expect(sim.Validate(bad)).To(MatchError("MeanInstancesPerApp must be 1 - 100"))
		})
	})
})
