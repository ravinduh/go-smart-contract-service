package datastruct

type Response struct {
	Data        *interface{}   `json:"data,omitempty"`
	CustomError *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	RequestId  string `json:"request_id,omitempty"`
	StatusCode int    `json:"-"` // internal service usage http.StatusBadRequest from net.http
	ErrorCode  int    `json:"error_code"`
	Message    string `json:"message,omitempty"`
}

type HealthCheckResponse struct {
	RequestId string `json:"request_id"`
	Message   string `json:"message"`
}
