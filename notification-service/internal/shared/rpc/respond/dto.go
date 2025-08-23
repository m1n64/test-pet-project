package respond

import "github.com/goccy/go-json"

const Version = "2.0"

type ID = any

type RPCErrorCode int

const (
	ParseError     RPCErrorCode = -32700
	InvalidRequest RPCErrorCode = -32600
	MethodNotFound RPCErrorCode = -32601
	InvalidParams  RPCErrorCode = -32602
	InternalError  RPCErrorCode = -32603
)

func (c RPCErrorCode) String() string {
	switch c {
	case ParseError:
		return "Parse error"
	case InvalidRequest:
		return "Invalid request"
	case MethodNotFound:
		return "Method not found"
	case InvalidParams:
		return "Invalid params"
	case InternalError:
		return "Internal error"
	default:
		return "Unknown"
	}
}

type Request struct {
	JSONRPC string          `json:"jsonrpc" binding:"required,eq=2.0"`
	Method  string          `json:"method"  binding:"required"`
	Params  json.RawMessage `json:"params,omitempty"` // by-name map[string]any by-pos []any
	ID      ID              `json:"id,omitempty"`
}

type Response[T any] struct {
	JSONRPC string    `json:"jsonrpc"`
	Result  *T        `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
	ID      ID        `json:"id"`
}

type RPCError struct {
	Code    RPCErrorCode `json:"code"`
	Message string       `json:"message"`
	Data    *ErrorDTO    `json:"data,omitempty"`
}

type ErrorDTO struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
