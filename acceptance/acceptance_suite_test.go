package acceptance_test

import (
	"math/rand"
	"net"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"

	"testing"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

const packagePath = "github.com/rosenhouse/cnsim"

var pathToServer string
var agoutiDriver *agouti.WebDriver

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error
	pathToServer, err = gexec.Build(packagePath)
	Expect(err).NotTo(HaveOccurred())
	return []byte(pathToServer)
}, func(crossNodeData []byte) {
	pathToServer = string(crossNodeData)
	rand.Seed(config.GinkgoConfig.RandomSeed + int64(GinkgoParallelNode()))

	agoutiDriver = agouti.PhantomJS()
	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = SynchronizedAfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())
}, func() {
	gexec.CleanupBuildArtifacts()
})

func VerifyTCPConnection(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
