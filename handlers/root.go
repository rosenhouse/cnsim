package handlers

import (
	"net/http"

	"code.cloudfoundry.org/lager"
)

type Root struct {
	Logger lager.Logger
}

func (h *Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger.Session("root")
	logger.Info("start")
	defer logger.Info("done")
	w.Write([]byte(body))
}

const body = `
<html>
  <head>
    <title>cnsim</title>
  </head>
  <body>
	  Inputs:
		<form action="/steady_state" method="GET">
		  <p> Hosts (1 - 1000): <input type="number" name="Hosts" min="1" max="1000"> </p>
		  <p> Apps (1 - 65k): <input type="number" name="Apps" min="1" max="65534"> </p>
		  <p> Instances / App (1 - 100): <input type="number" name="InstancesPerApp" min="1" max="100"> </p>
			<input type="submit" />
		</form>
  </body>
<html>
`
