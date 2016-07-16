package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/rosenhouse/cnsim/models"

	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter -o ../fakes/steady_state_simulator.go --fake-name SteadyStateSimulator . steadyStateSimulator
type steadyStateSimulator interface {
	Execute(logger lager.Logger, req models.SteadyStateRequest) (*models.SteadyStateResponse, error)
	Validate(req models.SteadyStateRequest) error
}

type SteadyState struct {
	Logger    lager.Logger
	Simulator steadyStateSimulator
}

func tryEncode(logger lager.Logger, w http.ResponseWriter, resp interface{}) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("encode-json", err)
	}
}

func (h *SteadyState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger.Session("steady-state")
	logger.Info("start")
	defer logger.Info("done")

	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		logger.Error("parse-form", err)
		w.WriteHeader(http.StatusBadRequest)

		tryEncode(logger, w, models.APIError{Error: fmt.Sprintf("parse-form: %s", err)})
		return
	}

	reqData := models.SteadyStateRequest{}
	decoder := schema.NewDecoder()
	err = decoder.Decode(&reqData, r.Form)
	if err != nil {
		logger.Error("decode", err)
		w.WriteHeader(http.StatusBadRequest)

		tryEncode(logger, w, models.APIError{Error: fmt.Sprintf("decode: %s", err)})
		return
	}

	err = h.Simulator.Validate(reqData)
	if err != nil {
		logger.Error("validation", err)
		w.WriteHeader(http.StatusBadRequest)
		tryEncode(logger, w, models.APIError{Error: fmt.Sprintf("validation: %s", err)})
		return
	}

	resp, err := h.Simulator.Execute(logger.Session("execute"), reqData)
	if err != nil {
		logger.Error("simulator", err)
		w.WriteHeader(http.StatusInternalServerError)

		tryEncode(logger, w, models.APIError{Error: fmt.Sprintf("simulator: %s", err)})
		return
	}

	tryEncode(logger, w, resp)
}
