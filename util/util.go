package util

import (
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	MaxIdleConnections    = 5
	MaxOpenConnections    = 10
	ConnectionLifeTimeInS = 60
	DateFormat            = "2006-01-02"
	DateTimeFormat        = "2006-01-02 15:04:05"
)

func GetEnvVariable(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func GetLogger(logLevel string) log.Logger {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.NewSyncLogger(logger)
	switch logLevel {
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	default:
		logger = level.NewFilter(logger, level.AllowDebug())
	}
	logger = log.With(logger,
		"Service", GetEnvVariable("SERVICE_NAME", "service"),
		"TimeStamp", log.DefaultTimestampUTC,
		"Caller", log.DefaultCaller,
	)
	return logger
}

func SetCorsHeaders(h http.Handler, logger log.Logger) http.HandlerFunc {
	// getting all the allowed origins into an array
	allowedOriginsSuffixesString := GetEnvVariable("CORS_ALLOWED_ORIGIN_SUFFIX", ".zoombookdirect.com,localhost:3000")
	allowedOriginsSuffixes := strings.Split(allowedOriginsSuffixesString, ",")

	fmt.Println("SetCorsHeaders SetCorsHeaders SetCorsHeaders SetCorsHeaders")
	return func(w http.ResponseWriter, r *http.Request) {

		originOfRequest := r.Header.Get("Origin")

		isAllowed := false
		for _, suffix := range allowedOriginsSuffixes {
			isAllowed = isAllowed || strings.HasSuffix(originOfRequest, suffix)
		}

		// check if request is allowed from the origin or not
		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", originOfRequest)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, x-csrf-token")
		}
		h.ServeHTTP(w, r)
	}
}

func GetSQLClient(logger log.Logger, maxIdleConns, maxOpenConns, connMaxLifeTimeInS int) (*sql.DB, error) {
	sqlDriverName := GetEnvVariable("SQL_DRIVER_NAME", "mysql")

	sqlDataSource := fmt.Sprintf("%s:%s@%s(%s)/%s",
		GetEnvVariable("DB_USERNAME", "userWalletService"),
		GetEnvVariable("DB_PASSWORD", "qazxsw123"),
		GetEnvVariable("DB_PROTOCOL", "tcp"),
		GetEnvVariable("DB_HOST", "db:3306"),
		GetEnvVariable("DB_NAME", "smart_contract_service"))

	sqlClient, err := sql.Open(sqlDriverName, sqlDataSource)
	if err != nil {
		_ = level.Error(logger).Log("Error", err)
		return nil, err
	}

	retries := 5
	for retries > 0 {
		err = sqlClient.Ping()
		if err == nil {
			_ = level.Debug(logger).Log("Message", "Connected to DB")
			break
		}

		_ = level.Error(logger).Log("Error", err)
		retries--
		time.Sleep(5 * time.Second) // Add a delay before retrying
	}
	if retries == 0 {
		_ = level.Error(logger).Log("Error", err)
		return nil, err
	}

	sqlClient.SetMaxIdleConns(maxIdleConns)
	sqlClient.SetMaxOpenConns(maxOpenConns)
	sqlClient.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifeTimeInS))
	return sqlClient, nil
}
