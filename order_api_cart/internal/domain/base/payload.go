package base

type ErrorInfo struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}
