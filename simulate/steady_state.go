package simulate

import (
	"fmt"

	"code.cloudfoundry.org/lager"

	"github.com/rosenhouse/cnsim/models"
)

type SteadyState struct{}

func (s *SteadyState) Execute(logger lager.Logger, req models.SteadyStateRequest) (*models.SteadyStateResponse, error) {
	logger.Info("start", lager.Data{"input": req})
	defer logger.Info("done")

	var resp models.SteadyStateResponse
	resp.Request = req
	totalInstances := float64(req.Apps) * float64(req.InstancesPerApp)
	resp.MeanInstancesPerHost = totalInstances / float64(req.Hosts)
	logger.Info("success", lager.Data{"output": resp})
	return &resp, nil
}

func validateRange(noun string, value int, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be %d - %d", noun, min, max)
	}
	return nil
}

func (s *SteadyState) Validate(req models.SteadyStateRequest) error {
	if err := validateRange("Hosts", req.Hosts, 1, 1000); err != nil {
		return err
	}
	if err := validateRange("Apps", req.Apps, 1, 65534); err != nil {
		return err
	}
	if err := validateRange("InstancesPerApp", req.InstancesPerApp, 1, 100); err != nil {
		return err
	}
	return nil
}
