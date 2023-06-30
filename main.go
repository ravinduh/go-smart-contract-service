package main

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"go-smart-contract-service/middleware"
	"go-smart-contract-service/service"
	"go-smart-contract-service/transport"
	"go-smart-contract-service/transport/routes"
	"go-smart-contract-service/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// for ENV variables' name always use CAPS with underscore
	httpAddr := util.GetEnvVariable("HTTP_ADDR", ":9124")
	logLevel := util.GetEnvVariable("LOG_LEVEL", "debug")

	logger := util.GetLogger(logLevel)

	svc := service.NewSmartContractService(logger)
	svc = middleware.NewSCSLoggingMiddleware(logger)(svc)

	endpoints := transport.MakeEndpoints(svc)

	// create a error channel, which can be used to stop the application in proper manner otherwise port will not get free in local
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
		_ = level.Info(logger).Log("Message", "stopping the Rate Service", "ErrChan", <-errChan)
	}()

	// Start the server listener
	go func() {
		h := routes.NewApi(endpoints, nil, logger)
		_ = level.Info(logger).Log("Transport", "HTTP", "Addr", httpAddr)
		server1 := &http.Server{
			Addr:    httpAddr,
			Handler: h,
		}

		errChan <- server1.ListenAndServe()
	}()

	_ = level.Error(logger).Log("Error", <-errChan)
}
