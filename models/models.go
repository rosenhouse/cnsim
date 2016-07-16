package models

type SteadyStateRequest struct {
	Hosts           int
	Apps            int
	InstancesPerApp int
}

type SteadyStateResponse struct {
	SteadyStateRequest

	MeanInstancesPerHost float64
}

type APIError struct {
	Error string
}
