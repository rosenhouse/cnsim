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
	Id   int `json:"-"`
	Size int `json:"s"`
}

type Instance struct {
	Id     int `json:"-"`
	AppId  int `json:"a"`
	HostId int `json:"h"`
}

type APIError struct {
	Error string
}
