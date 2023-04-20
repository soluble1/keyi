package recovery

import (
	"github.com/soluble1/mweb"
)

type MiddlewareHandler struct {
	code    int
	errMsg  []byte
	rcvFunc func(ctx *mweb.Context)
}

func NewRecoverHandle(code int, errMsg []byte, fun func(ctx *mweb.Context)) *MiddlewareHandler {
	return &MiddlewareHandler{
		code:    code,
		errMsg:  errMsg,
		rcvFunc: fun,
	}
}

func (m *MiddlewareHandler) Build() mweb.Middleware {
	return func(next mweb.HandlerFunc) mweb.HandlerFunc {
		return func(ctx *mweb.Context) {
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
