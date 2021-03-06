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
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
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
			NumHosts:            1000,
			NumApps:             10000,
			MeanInstancesPerApp: 2,
		}
		var responseData models.SteadyStateResponse
		var apiError models.APIError
		resp, err := apiClient.New().Get("/steady_state").QueryStruct(requestData).Receive(&responseData, &apiError)
		Expect(err).NotTo(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(200))
		By("checking the CORS response")
		Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))

		By("checking the original request is included with the response")
		Expect(responseData.Request).To(Equal(requestData))

		By("checking the mean instances per host")
		Expect(responseData.MeanInstancesPerHost).To(Equal(20.0))

		By("checking the mean instances per host")
		Expect(responseData.Apps).To(HaveLen(10000))
	})

	Describe("Web form", func() {
		var page *agouti.Page

		BeforeEach(func() {
			var err error
			page, err = agoutiDriver.NewPage()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(page.Destroy()).To(Succeed())
		})
		It("provide a web interface to the API", func() {
			By("having a root page", func() {
				Expect(page.Navigate("http://" + address)).To(Succeed())
			})

			By("allowing the user to fill out the form and run a simulation", func() {
				Eventually(page.FindByName("NumHosts")).Should(BeFound())
				Expect(page.FindByName("NumHosts").Fill("100")).To(Succeed())
				Expect(page.FindByName("NumApps").Fill("300")).To(Succeed())
				Expect(page.FindByName("MeanInstancesPerApp").Fill("5")).To(Succeed())
				Expect(page.FindByButton("Simulate").Click()).To(Succeed())
			})

			By("showing a histogram of app sizes", func() {
				Eventually(page.HTML).Should(ContainSubstring(`Size (instances)`))
				Eventually(page.HTML).Should(ContainSubstring(`40</text>`)) // tick mark on y axis
			})
		})
	})

})
