package respond

import (
	"net/http"
)

type Builder[T any] struct {
	id   ID
	resp Response[T]
}

func New[T any](c Ctx, id ID) *Builder[T] {
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
	b.resp.Error = &RPCError{
		Code:    rpcCode,
		Message: message,
		Data: &ErrorDTO{
			Code:    appCode,
			Message: message,
			Details: details,
		},
	}

	return b
}

func (b *Builder[T]) JSON(c Ctx) {
	if b.id == nil {
		return
	}

	c.JSON(http.StatusOK, b.resp)

}

func OK[T any](c Ctx, id ID, v T) {
	New[T](c, id).Data(v).JSON(c)
}

func Fail(c Ctx, id ID, rpcCode RPCErrorCode, appCode, message string, details interface{}) {
	New[struct{}](c, id).Error(rpcCode, appCode, message, details).JSON(c)
}
