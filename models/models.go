package models

type SteadyStateRequest struct {
	NumHosts            int
	NumApps             int
	MeanInstancesPerApp int
}

type SteadyStateResponse struct {
	Request SteadyStateRequest

	MeanInstancesPerHost float64
	Apps                 []App
}

type App struct {
	Id   int
	Size int
}

type APIError struct {
	Error string
}
