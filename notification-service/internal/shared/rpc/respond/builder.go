package respond

import (
	"net/http"
)

type Builder[T any] struct {
	id   ID
	resp Response[T]
}

func New[T any](id ID) *Builder[T] {
	return &Builder[T]{
		id: id,
		resp: Response[T]{
			JSONRPC: Version,
			ID:      id,
		},
	}
}

func (b *Builder[T]) Data(v T) *Builder[T] {
	b.resp.Result = &v
	b.resp.Error = nil
	return b
}

func (b *Builder[T]) Error(rpcCode RPCErrorCode, appCode, message string, details interface{}) *Builder[T] {
	b.resp.Result = nil
	b.resp.Error = NewRPCError(rpcCode, appCode, message, details)

	return b
}

func NewRPCError(rpcCode RPCErrorCode, appCode, message string, details any) *RPCError {
	return &RPCError{
		Code:    rpcCode,
		Message: message,
		Data: &ErrorDTO{
			Code:    appCode,
			Message: message,
			Details: details,
		},
	}
}

func (b *Builder[T]) JSON(c Ctx) {
	c.JSON(http.StatusOK, b.resp)
}

func (b *Builder[T]) GetResponse() Response[T] {
	return b.resp
}

func BuildOK[T any](id ID, v T) Response[T] {
	return New[T](id).Data(v).GetResponse()
}

func BuildFail(id ID, rpcCode RPCErrorCode, appCode, message string, details interface{}) Response[any] {
	return New[any](id).Error(rpcCode, appCode, message, details).GetResponse()
}

func OK[T any](c Ctx, id ID, v T) {
	New[T](id).Data(v).JSON(c)
}

func Fail(c Ctx, id ID, rpcCode RPCErrorCode, appCode, message string, details interface{}) {
	New[struct{}](id).Error(rpcCode, appCode, message, details).JSON(c)
}
