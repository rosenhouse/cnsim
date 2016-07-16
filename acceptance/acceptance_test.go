package acceptance_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"

	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/rosenhouse/cnsim/models"
)

var _ = Describe("CNSim Server", func() {
	var (
		session *gexec.Session
		address string

		apiClient *sling.Sling
	)

	var serverIsAvailable = func() error {
		return VerifyTCPConnection(address)
	}

	BeforeEach(func() {
		port := 10000 + rand.Intn(10000)
		serverCmd := exec.Command(pathToServer)
		serverCmd.Env = []string{fmt.Sprintf("PORT=%d", port)}
		var err error
		session, err = gexec.Start(serverCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		address = fmt.Sprintf("127.0.0.1:%d", port)
		apiClient = sling.New().Base("http://" + address).Client(http.DefaultClient)

		Eventually(serverIsAvailable, "5s").Should(Succeed())
	})

	AfterEach(func() {
		session.Interrupt()
		Eventually(session, "5s").Should(gexec.Exit())
	})

	It("should say hello on /", func() {
		resp, err := http.Get("http://" + address + "/")
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(bodyBytes).To(ContainSubstring("<html>"))
	})

	It("should compute steady-state stats on /steady_state", func() {
		requestData := models.SteadyStateRequest{
			Hosts:           1000,
			Apps:            10000,
			InstancesPerApp: 2,
		}
		var responseData models.SteadyStateResponse
		var apiError models.APIError
		resp, err := apiClient.New().Get("/steady_state").QueryStruct(requestData).Receive(&responseData, &apiError)
		Expect(err).NotTo(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(200))

		By("checking the original request is included with the response")
		Expect(responseData.Request).To(Equal(requestData))

		By("checking the mean instances per host")
		Expect(responseData.MeanInstancesPerHost).To(Equal(20.0))
	})
})
