package repository

import (
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go-smart-contract-service/util"
	"time"
)

type UserWalletRepo struct {
	logger    log.Logger
	sqlClient *sql.DB
}

func (uw UserWalletRepo) InsertUserRecord(requestId, nric, encryptedWalletAddress string) error {
	currentTime := time.Now().UTC().Format(util.DateTimeFormat)
	_, err := uw.sqlClient.Exec("INSERT INTO user_wallet_log (nric, wallet_address, created_at) VALUES (?, ?, ?)",
		nric, encryptedWalletAddress, currentTime,
	)
	if err != nil {
		_ = level.Error(uw.logger).Log("RequestId", requestId, "Error", err)
		return err
	}
	return nil
}

func (uw UserWalletRepo) GetUserWalletCount(requestId, nric string) (int, error) {
	var walletCount int
	err := uw.sqlClient.QueryRow(
		"SELECT count(1) FROM user_wallet_log WHERE nric = ? AND deleted_at IS NULL", nric).Scan(&walletCount)
	if err != nil {
		_ = level.Error(uw.logger).Log("RequestId", requestId, "Error", err)
		return walletCount, err
	}
	return walletCount, nil
}

func (uw UserWalletRepo) GetUserWallet(requestId, nric string) (*string, error) {

	var encryptedWalletAddress *string
	err := uw.sqlClient.QueryRow(
		"SELECT wallet_address FROM user_wallet_log WHERE nric = ? AND deleted_at IS NULL", nric).Scan(&encryptedWalletAddress)
	if err != nil {
		_ = level.Error(uw.logger).Log("RequestId", requestId, "Error", err)
		return nil, err
	}
	return encryptedWalletAddress, nil
}

func NewUserWalletRepo(logger log.Logger, sqlClient *sql.DB) UserWalletRepo {
	return UserWalletRepo{
		logger:    logger,
		sqlClient: sqlClient,
	}
}
