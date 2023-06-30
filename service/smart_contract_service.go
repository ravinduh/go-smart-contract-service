package service

import (
	"github.com/go-kit/kit/log"
	"go-smart-contract-service/datastruct"
)

type SmartContractService interface {
	HealthCheck(requestId string) datastruct.Response
	Test(requestId string) datastruct.Response
}

type smartContractService struct {
	logger log.Logger
}

func (scs smartContractService) HealthCheck(requestId string) datastruct.Response {
	var res interface{}
	res = datastruct.HealthCheckResponse{
		RequestId: requestId,
		Message:   "Service is up",
	}
	return datastruct.Response{
		Data: &res,
	}
}

func (scs smartContractService) Test(requestId string) datastruct.Response {
	return datastruct.Response{}
}

// NewSmartContractService returns a naive, stateless implementation of SmartContractService.
func NewSmartContractService(logger log.Logger) SmartContractService {
	return &smartContractService{
		logger: logger,
	}
}
