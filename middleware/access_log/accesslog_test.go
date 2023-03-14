package access_log

import (
	"testing"
	"web_copy"
)

func TestAccessLog(t *testing.T) {
	server := web_copy.NewHTTPServer()

	acclogMid := NewAccessLog()

	server.Use(acclogMid.Build())

	server.Get("/", func(ctx *web_copy.Context) {
		ctx.Resp.Write([]byte("this is /"))
	})

	server.Get("/user/name/age", func(ctx *web_copy.Context) {
		ctx.Resp.Write([]byte("this is /user/name/age"))
	})

	server.Get("/user/name", func(ctx *web_copy.Context) {
		ctx.Resp.Write([]byte("this is /user/name"))
	})

	server.Start(":8081")
}