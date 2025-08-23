package respond

type Response[T any] struct {
	Success bool      `json:"success"`
	Data    *T        `json:"data,omitempty"`
	Error   *ErrorDTO `json:"error,omitempty"`
	//Meta    MetaDTO   `json:"meta"`
}

type MetaDTO struct {
	RequestID string `json:"request_id,omitempty"`
}

type ErrorDTO struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
