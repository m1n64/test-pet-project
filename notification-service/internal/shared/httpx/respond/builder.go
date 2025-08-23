package respond

import (
	"net/http"
)

type Builder[T any] struct {
	status int
	resp   Response[T]
}

func New[T any](c Ctx) *Builder[T] {
	return &Builder[T]{
		status: http.StatusOK,
		resp: Response[T]{
			Success: true,
			/*Meta: MetaDTO{
				RequestID: rid,
			},*/
		},
	}
}

func (b *Builder[T]) Status(s int) *Builder[T] {
	b.status = s
	return b
}

func (b *Builder[T]) Data(v T) *Builder[T] {
	b.resp.Data = &v
	b.resp.Success = true
	b.resp.Error = nil
	return b
}

func (b *Builder[T]) Error(code, message string, details interface{}) *Builder[T] {
	b.resp.Success = false
	b.resp.Error = &ErrorDTO{
		Code:    code,
		Message: message,
		Details: details,
	}
	b.resp.Data = nil
	return b
}

func (b *Builder[T]) JSON(c Ctx) {
	c.JSON(b.status, b.resp)
}

func OK[T any](c Ctx, v T) {
	New[T](c).Data(v).JSON(c)
}

func Fail(c Ctx, status int, code, message string, details interface{}) {
	New[struct{}](c).Status(status).Error(code, message, details).JSON(c)
}
