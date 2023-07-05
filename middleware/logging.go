package middleware

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go-smart-contract-service/datastruct"
	"go-smart-contract-service/service"
	"time"
)

type SCSLoggingMiddleware func(svc service.SmartContractService) service.SmartContractService

type scsLogging struct {
	logger log.Logger
	next   service.SmartContractService
}

func (scsl scsLogging) HealthCheck(requestId string) datastruct.Response {
	defer func(begin time.Time) {
		took := int64(time.Since(begin) / time.Millisecond)
		_ = level.Info(scsl.logger).Log(
			"Method", "HealthCheck",
			"RequestId", requestId,
			"UserId", 0,
			"Request", nil,
			"ValErr", nil,
			"Took", took,
		)
	}(time.Now())

	return scsl.next.HealthCheck(requestId)
}

func (scsl scsLogging) GetReceipt(requestId string, receiptReq datastruct.ReceiptRequest) datastruct.Response {
	defer func(begin time.Time) {
		took := int64(time.Since(begin) / time.Millisecond)
		_ = level.Info(scsl.logger).Log(
			"Method", "GetReceipt",
			"RequestId", requestId,
			"UserId", 0,
			"Request", receiptReq.String(),
			"ValErr", nil,
			"Took", took,
		)
	}(time.Now())

	return scsl.next.GetReceipt(requestId, receiptReq)
}

func (scsl scsLogging) Test(requestId string) datastruct.Response {
	var valErr error
	defer func(begin time.Time) {
		took := int64(time.Since(begin) / time.Millisecond)
		_ = level.Info(scsl.logger).Log(
			"Method", "Test",
			"RequestId", requestId,
			"UserId", 0,
			"Request", "",
			"ValErr", valErr,
			"Took", took,
		)
	}(time.Now())

	return scsl.next.Test(requestId)
}

// NewSCSLoggingMiddleware returns a naive, stateless implementation of SCSLoggingMiddleware.
func NewSCSLoggingMiddleware(logger log.Logger) SCSLoggingMiddleware {
	return func(next service.SmartContractService) service.SmartContractService {
		return &scsLogging{
			logger: logger,
			next:   next,
		}
	}
}
