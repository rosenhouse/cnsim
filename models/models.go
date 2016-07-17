package models

type SteadyStateRequest struct {
	NumHosts            int
	NumApps             int
	MeanInstancesPerApp int
}

type SteadyStateResponse struct {
	Request SteadyStateRequest

	MeanInstancesPerHost float64
	TotalInstances       int
	Apps                 []App
	Instances            []Instance
}

type App struct {
	Id   int
	Size int
}

type Instance struct {
	Id     int
	AppId  int
	HostId int
}

type APIError struct {
	Error string
}
