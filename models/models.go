package models

type SteadyStateRequest struct {
	Hosts           int
	Apps            int
	InstancesPerApp int
}

type SteadyStateResponse struct {
	Request SteadyStateRequest

	MeanInstancesPerHost float64
}

type APIError struct {
	Error string
}
