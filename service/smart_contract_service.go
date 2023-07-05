package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go-smart-contract-service/datastruct"
	"go-smart-contract-service/repository"
)

type SmartContractService interface {
	HealthCheck(requestId string) datastruct.Response
	GetReceipt(requestId string, receiptRequest datastruct.ReceiptRequest) datastruct.Response
	Test(requestId string) datastruct.Response
}

type smartContractService struct {
	logger           log.Logger
	encryptionHelper EncryptionHelper
	userWalletRepo   repository.UserWalletRepo
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

func (scs smartContractService) GetReceipt(requestId string, receiptRequest datastruct.ReceiptRequest) datastruct.Response {
	var res interface{}

	walletCount, err := scs.userWalletRepo.GetUserWalletCount(requestId, receiptRequest.NRIC)
	if err != nil {
		_ = level.Error(scs.logger).Log("RequestId", requestId, "Error", err)
		return datastruct.Response{
			CustomError: &datastruct.ErrorResponse{
				RequestId:  requestId,
				StatusCode: 500,
				ErrorCode:  1001,
				Message:    err.Error(),
			},
		}
	}

	// Check if nric already associated with a wallet, if yes then return error
	if walletCount > 0 {
		encryptedUserWallet, err := scs.userWalletRepo.GetUserWallet(requestId, receiptRequest.NRIC)
		if err != nil {
			_ = level.Error(scs.logger).Log("RequestId", requestId, "Error", err)
			return datastruct.Response{
				CustomError: &datastruct.ErrorResponse{
					RequestId:  requestId,
					StatusCode: 500,
					ErrorCode:  1002,
					Message:    err.Error(),
				},
			}
		}

		userWalletStored, err := scs.encryptionHelper.DecryptBase64(requestId, *encryptedUserWallet)
		if err != nil {
			_ = level.Error(scs.logger).Log("RequestId", requestId, "Error", err)
			return datastruct.Response{
				CustomError: &datastruct.ErrorResponse{
					RequestId:  requestId,
					StatusCode: 500,
					ErrorCode:  1003,
					Message:    err.Error(),
				},
			}
		}

		if userWalletStored != receiptRequest.WalletAddress {
			return datastruct.Response{
				CustomError: &datastruct.ErrorResponse{
					RequestId:  requestId,
					StatusCode: 500,
					ErrorCode:  1004,
					Message:    "NRIC is already associated with another wallet",
				},
			}
		}

	} else {
		// if the nric is new then insert and return
		encryptedNewWallet, err := scs.encryptionHelper.EncryptBase64(requestId, receiptRequest.WalletAddress)
		err = scs.userWalletRepo.InsertUserRecord(requestId, receiptRequest.NRIC, string(encryptedNewWallet))
		if err != nil {
			_ = level.Error(scs.logger).Log("RequestId", requestId, "Error", err)
			return datastruct.Response{
				CustomError: &datastruct.ErrorResponse{
					RequestId:  requestId,
					StatusCode: 500,
					ErrorCode:  1005,
					Message:    "Error while inserting user record",
				},
			}
		}
	}

	hash, err := getHash(requestId, scs.logger, receiptRequest)
	if err != nil {
		_ = level.Error(scs.logger).Log("RequestId", requestId, "Error", err)
		return datastruct.Response{
			CustomError: &datastruct.ErrorResponse{
				RequestId:  requestId,
				StatusCode: 500,
				ErrorCode:  1006,
				Message:    "Error while hashing the data",
			},
		}
	}

	res = datastruct.ReceiptResponse{
		RequestId: requestId,
		Receipt:   *hash,
	}
	return datastruct.Response{
		Data: &res,
	}
}

func getHash(requestId string, logger log.Logger, request datastruct.ReceiptRequest) (*string, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		_ = level.Error(logger).Log("RequestId", requestId, "Error", err)
		return nil, err
	}

	// Hash the JSON data using SHA-256
	hash := sha256.Sum256(jsonData)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash[:])

	return &hashString, nil
}

func (scs smartContractService) Test(requestId string) datastruct.Response {
	return datastruct.Response{}
}

// NewSmartContractService returns a naive, stateless implementation of SmartContractService.
func NewSmartContractService(logger log.Logger, encryptionHelper EncryptionHelper,
	userWalletRepo repository.UserWalletRepo) SmartContractService {
	return &smartContractService{
		logger:           logger,
		encryptionHelper: encryptionHelper,
		userWalletRepo:   userWalletRepo,
	}
}
