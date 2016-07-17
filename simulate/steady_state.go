package simulate

import (
	"fmt"

	"code.cloudfoundry.org/lager"

	"github.com/rosenhouse/cnsim/models"
)

//go:generate counterfeiter -o ../fakes/mean_parameterized_discrete_distribution.go --fake-name MeanParameterizedDiscreteDistribution . meanParameterizedDiscreteDistribution
type meanParameterizedDiscreteDistribution interface {
	Sample(mean float64) (int, error)
}

type SteadyState struct {
	AppSizeDistribution meanParameterizedDiscreteDistribution
}

func (s *SteadyState) Execute(logger lager.Logger, req models.SteadyStateRequest) (*models.SteadyStateResponse, error) {
	logger.Info("start", lager.Data{"input": req})
	defer logger.Info("done")

	var resp models.SteadyStateResponse
	resp.Request = req
	totalInstances := float64(req.NumApps) * float64(req.MeanInstancesPerApp)
	resp.MeanInstancesPerHost = totalInstances / float64(req.NumHosts)
	resp.Apps = make([]models.App, req.NumApps)
	var err error
	for i, _ := range resp.Apps {
		resp.Apps[i].Id = i
		resp.Apps[i].Size, err = s.AppSizeDistribution.Sample(float64(req.MeanInstancesPerApp))
		if err != nil {
			return nil, fmt.Errorf("sampling app size: %s", err)
		}
	}

	logger.Info("success")
	return &resp, nil
}

func validateRange(noun string, value int, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be %d - %d", noun, min, max)
	}
	return nil
}

func (s *SteadyState) Validate(req models.SteadyStateRequest) error {
	if err := validateRange("NumHosts", req.NumHosts, 1, 1000); err != nil {
		return err
	}
	if err := validateRange("NumApps", req.NumApps, 1, 65534); err != nil {
		return err
	}
	if err := validateRange("MeanInstancesPerApp", req.MeanInstancesPerApp, 1, 100); err != nil {
		return err
	}
	return nil
}
