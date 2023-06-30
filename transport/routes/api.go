package routes

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go-smart-contract-service/datastruct"
	"go-smart-contract-service/transport"
	"go-smart-contract-service/util"
	"net/http"
)

// NewApi wires Go kit endpoints to the HTTP transport.
func NewApi(svcEndpoints transport.Endpoints, options []kitHttp.ServerOption, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// handle CORS
	r.Methods(http.MethodOptions).Handler(util.SetCorsHeaders(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			return
		}), logger))

	// HTTP GET - /service/health-check
	r.Methods(http.MethodGet).Path("/health-check").Handler(kitHttp.NewServer(
		svcEndpoints.HealthCheck,
		decodeRequest,
		encodeResponse,
	))

	// HTTP Post - /service/test
	r.Methods(http.MethodPost).Path("/service/test").Handler(kitHttp.NewServer(
		svcEndpoints.Test,
		decodeRequest,
		encodeResponse,
	))

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res, _ := response.(datastruct.Response)
	w.Header().Add("Content-Type", "application/json")
	if res.CustomError == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(res.CustomError.StatusCode)
	}

	return json.NewEncoder(w).Encode(response)
}
