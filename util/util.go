package util

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	"strings"
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
