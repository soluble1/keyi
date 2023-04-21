package cors

import (
	"github.com/soluble1/mweb"
	"net/http"
)

type MiddlewareBuilder struct {
	AllowOrigin string
}

func (m MiddlewareBuilder) Build() mweb.Middleware {
	return func(next mweb.HandlerFunc) mweb.HandlerFunc {
		return func(ctx *mweb.Context) {
			allowOrigin := m.AllowOrigin
			if allowOrigin == "" {
				allowOrigin = ctx.Req.Header.Get("Origin")
			}
			ctx.Resp.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			// ctx.Resp.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Resp.Header().Set("Access-Control-Allow-Credentials", "true")
			if ctx.Resp.Header().Get("Access-Control-Allow-Headers") == "" {
				ctx.Resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
			}
			if ctx.Req.Method == http.MethodOptions {
				ctx.RespStatusCode = 200
				ctx.RespData = []byte("ok")
				return
			}
			next(ctx)
		}
	}
}
