package main

import (
	"fmt"
	"log"
	"os"

	"code.cloudfoundry.org/lager"

	"github.com/rosenhouse/cnsim/handlers"
	"github.com/rosenhouse/cnsim/simulate"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
	"github.com/tedsuo/rata"
)

func getEnv(logger lager.Logger, name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
		logger.Info("missing-env-var", lager.Data{"name": name, "default-to": value})
	} else {
		logger.Info("read-env-var", lager.Data{"name": name, "value": value})
	}
	return value
}

func main() {
	logger := lager.NewLogger("cnsim-server")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	port := getEnv(logger, "PORT", "9000")
	listenAddress := getEnv(logger, "LISTEN_ADDRESS", "127.0.0.1")

	address := fmt.Sprintf("%s:%s", listenAddress, port)
	logger.Info("listen", lager.Data{"address": address})

	routes := rata.Routes{
		{Name: "root", Method: "GET", Path: "/"},
		{Name: "steady_state", Method: "GET", Path: "/steady_state"},
	}

	rataHandlers := rata.Handlers{
		"root": &handlers.Root{
			Logger: logger,
		},
		"steady_state": &handlers.SteadyState{
			Logger:    logger,
			Simulator: &simulate.SteadyState{},
		},
	}

	router, err := rata.NewRouter(routes, rataHandlers)
	if err != nil {
		log.Fatalf("unable to create rata Router: %s", err) // not tested
	}

	monitor := ifrit.Invoke(sigmon.New(grouper.NewOrdered(os.Interrupt, grouper.Members{
		{"http_server", http_server.New(address, router)},
	})))
	err = <-monitor.Wait()
	if err != nil {
		log.Fatalf("ifrit: %s", err)
	}
}
