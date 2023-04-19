package err_handle

import "github.com/soluble1/mweb"

type MiddlewareHandler struct {
	errPage map[int][]byte
}

func NewErrHandle() *MiddlewareHandler {
	return &MiddlewareHandler{
		errPage: make(map[int][]byte, 64),
	}
}

func (m *MiddlewareHandler) AddErrPage(code int, data []byte) {
	m.errPage[code] = data
}

func (m *MiddlewareHandler) Build() web_copy.Middleware {
	return func(next web_copy.HandlerFunc) web_copy.HandlerFunc {
		return func(ctx *web_copy.Context) {
			defer func() {
				page, ok := m.errPage[ctx.RespStatusCode]
				if ok {
					ctx.RespData = page
				}
			}()

			next(ctx)
		}
	}
}
