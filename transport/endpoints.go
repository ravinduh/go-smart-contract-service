package transport

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"go-smart-contract-service/datastruct"
	"go-smart-contract-service/service"
	"net/http"
)

// Endpoints holds all Go kit endpoints for the  service.
type Endpoints struct {
	HealthCheck endpoint.Endpoint
	Test        endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the  service.
func MakeEndpoints(svc service.SmartContractService) Endpoints {
	return Endpoints{
		HealthCheck: MakeHealthCheckEndpoint(svc),
		Test:        MakeTestEndpoint(svc),
	}
}

func MakeHealthCheckEndpoint(svc service.SmartContractService) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(*http.Request)
		requestId := getRequestIdFromHeader(req)
		return svc.HealthCheck(requestId), nil
	}
}

func MakeTestEndpoint(svc service.SmartContractService) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(*http.Request)
		requestId := getRequestIdFromHeader(req)
		var request datastruct.Request
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			return datastruct.Response{
				Data: nil,
				CustomError: &datastruct.ErrorResponse{
					RequestId:  requestId,
					StatusCode: http.StatusBadRequest,
					ErrorCode:  0,
					Message:    "",
				},
			}, err
		}
		return svc.Test(requestId), nil
	}
}

const RequestIdHeader = "Request-Id"

func getRequestIdFromHeader(req *http.Request) string {
	requestId := req.Header.Get(RequestIdHeader)
	if requestId == "" {
		requestId = uuid.New().String()
		req.Header.Add(RequestIdHeader, requestId)
	}
	return requestId
}
