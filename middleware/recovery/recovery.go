package recovery

import (
	"web_copy"
)

type MiddlewareHandler struct {
	code    int
	errMsg  []byte
	rcvFunc func(ctx *web_copy.Context)
}

func NewRecoverHandle(code int, errMsg []byte, fun func(ctx *web_copy.Context)) *MiddlewareHandler {
	return &MiddlewareHandler{
		code:    code,
		errMsg:  errMsg,
		rcvFunc: fun,
	}
}

func (m *MiddlewareHandler) Build() web_copy.Middleware {
	return func(next web_copy.HandlerFunc) web_copy.HandlerFunc {
		return func(ctx *web_copy.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.code
					ctx.RespData = m.errMsg
					m.rcvFunc(ctx)
				}
			}()

			next(ctx)
		}
	}
}
