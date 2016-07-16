package models

type SteadyStateRequest struct {
	NumHosts            int
	NumApps             int
	MeanInstancesPerApp int
}

type SteadyStateResponse struct {
	Request SteadyStateRequest

	MeanInstancesPerHost float64
}

type APIError struct {
	Error string
}
