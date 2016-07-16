package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/rosenhouse/cnsim/fakes"
	"github.com/rosenhouse/cnsim/handlers"
	"github.com/rosenhouse/cnsim/models"
)

var _ = Describe("Steady State Handler", func() {
	var (
		logger    *lagertest.TestLogger
		simulator *fakes.SteadyStateSimulator
		handler   handlers.SteadyState

		request   *http.Request
		response  *httptest.ResponseRecorder
		apiClient *sling.Sling
		reqData   models.SteadyStateRequest
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		simulator = &fakes.SteadyStateSimulator{}
		handler = handlers.SteadyState{
			Logger:    logger,
			Simulator: simulator,
		}

		response = httptest.NewRecorder()

		apiClient = sling.New().Base("http://localhost/").Client(http.DefaultClient)

		reqData = models.SteadyStateRequest{
			Hosts:           123,
			Apps:            456,
			InstancesPerApp: 789,
		}

		var err error
		request, err = apiClient.New().Get("steady_state").QueryStruct(reqData).Request()
		Expect(err).NotTo(HaveOccurred())

		simulator.ExecuteReturns(&models.SteadyStateResponse{
			SteadyStateRequest: reqData,

			MeanInstancesPerHost: 3.14159,
		}, nil)
	})

	It("unmarshals the request query data and validates it", func() {
		handler.ServeHTTP(response, request)

		Expect(simulator.ValidateCallCount()).To(Equal(1))
		Expect(simulator.ValidateArgsForCall(0)).To(Equal(reqData))
	})

	It("passes the data to the steady state simulator", func() {
		handler.ServeHTTP(response, request)

		Expect(simulator.ExecuteCallCount()).To(Equal(1))
		l, r := simulator.ExecuteArgsForCall(0)
		Expect(l.SessionName()).To(Equal("test.steady-state.execute"))
		Expect(r).To(Equal(reqData))
	})

	It("marshals the simulator response to JSON", func() {
		handler.ServeHTTP(response, request)

		Expect(response.Code).To(Equal(200))
		Expect(response.HeaderMap.Get("Content-Type")).To(Equal("application/json"))

		var respData models.SteadyStateResponse
		json.Unmarshal(response.Body.Bytes(), &respData)
		Expect(respData.MeanInstancesPerHost).To(Equal(3.14159))
	})

	Context("when parsing the form data fails", func() {
		BeforeEach(func() {
			request.URL.RawQuery = "%%%"
			handler.ServeHTTP(response, request)
		})

		It("logs the error", func() {
			Expect(logger.Buffer()).To(gbytes.Say(`"error":"invalid URL escape`))
		})

		It("responds with code 400 and the message in json", func() {
			Expect(response.Code).To(Equal(400))
			Expect(response.HeaderMap.Get("Content-Type")).To(Equal("application/json"))

			var err models.APIError
			Expect(json.Unmarshal(response.Body.Bytes(), &err)).To(Succeed())
			Expect(err.Error).To(HavePrefix("parse-form: invalid URL escape"))
		})
	})

	Context("when decoding the form data into a structure fails", func() {
		BeforeEach(func() {
			request.URL.RawQuery = "???"
			handler.ServeHTTP(response, request)
		})

		It("logs the error", func() {
			Expect(logger.Buffer()).To(gbytes.Say(`schema: invalid path`))
		})

		It("responds with code 400 and the message in json", func() {
			Expect(response.Code).To(Equal(400))
			Expect(response.HeaderMap.Get("Content-Type")).To(Equal("application/json"))

			var err models.APIError
			Expect(json.Unmarshal(response.Body.Bytes(), &err)).To(Succeed())
			Expect(err.Error).To(ContainSubstring("decode: schema: invalid path"))
		})
	})

	Context("when validating the data fails", func() {
		BeforeEach(func() {
			simulator.ValidateReturns(errors.New("banana"))
			handler.ServeHTTP(response, request)
		})

		It("logs the error", func() {
			Expect(logger.Buffer()).To(gbytes.Say(`banana`))
		})

		It("returns a useful error", func() {
			Expect(response.Code).To(Equal(400))
			Expect(response.HeaderMap.Get("Content-Type")).To(Equal("application/json"))

			var err models.APIError
			Expect(json.Unmarshal(response.Body.Bytes(), &err)).To(Succeed())
			Expect(err.Error).To(ContainSubstring("validation: banana"))
		})
	})

	Context("when the simulator errors", func() {
		BeforeEach(func() {
			simulator.ExecuteReturns(nil, errors.New("banana"))
			handler.ServeHTTP(response, request)
		})

		It("logs the error", func() {
			Expect(logger.Buffer()).To(gbytes.Say(`banana`))
		})

		It("responds with code 500 and the message in JSON", func() {
			Expect(response.Code).To(Equal(500))
			Expect(response.HeaderMap.Get("Content-Type")).To(Equal("application/json"))

			var err models.APIError
			Expect(json.Unmarshal(response.Body.Bytes(), &err)).To(Succeed())
			Expect(err.Error).To(ContainSubstring("simulator: banana"))
		})
	})

	Context("when writing the json response fails", func() {
		BeforeEach(func() {
			resp := &fakes.ResponseWriter{}
			resp.HeaderReturns(make(http.Header))
			resp.WriteReturns(0, errors.New("potato"))
			handler.ServeHTTP(resp, request)
		})

		It("logs the error", func() {
			Expect(logger.Buffer()).To(gbytes.Say(`potato`))
		})
	})
})
