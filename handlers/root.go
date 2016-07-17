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
<!DOCTYPE html>
<html lang="en">
	<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
			<title>cnsim</title>
			<meta charset="UTF-8">
			<link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.6/css/bootstrap.min.css">
			<link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/dc/1.7.5/dc.css">
	</head>
  <body>
		<form action="/steady_state" method="GET" id="steady-state-request">
		  <p> Num Hosts (1 - 1000): <input type="number" name="NumHosts" min="1" max="1000"> </p>
		  <p> Num Apps (1 - 65k): <input type="number" name="NumApps" min="1" max="65534"> </p>
		  <p> Avg Instances / App (1 - 100): <input type="number" name="MeanInstancesPerApp" min="1" max="100"> </p>
		</form>
		<button id="submit-button">Simulate</button>
		<div class="container">
			<div id="apps"></div>
		</div>
		<script type="text/javascript" src="//cdnjs.cloudflare.com/ajax/libs/jquery/3.1.0/jquery.slim.min.js"></script>
		<script type="text/javascript" src="//cdnjs.cloudflare.com/ajax/libs/d3/3.5.17/d3.min.js"></script>
		<script type="text/javascript" src="//cdnjs.cloudflare.com/ajax/libs/crossfilter/1.3.12/crossfilter.min.js"></script>
		<script type="text/javascript" src="//cdnjs.cloudflare.com/ajax/libs/dc/1.7.5/dc.min.js"></script>
		<script type="text/javascript">
			document.getElementById("submit-button").addEventListener("click", function () {
				var chartApps = dc.barChart("#apps");
				var jsonURL = "/steady_state?" + $('#steady-state-request').serialize()
				console.log(jsonURL)
				d3.json(jsonURL,
				 function(error, steady_state) {
					var apps = steady_state.Apps
					var ndx            = crossfilter(apps),
							countDimension = ndx.dimension(function(d) {return d.DesiredInstanceCount;}),
							appGroup       = countDimension.group().reduceCount();
					chartApps
						.width(768)
						.height(480)
						.x(d3.scale.linear().domain([1,25]))
						.brushOn(false)
						.yAxisLabel("# Apps")

						.xAxisLabel("Size (instances)")
						.dimension(countDimension)
						.group(appGroup);
						chartApps.xAxis().tickValues(d3.range(1,25,1));
						chartApps.render();
				});
			});
		</script>
  </body>
<html>
`
