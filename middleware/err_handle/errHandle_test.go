package err_handle

import (
	"testing"
	"web_copy"
	"web_copy/middleware/access_log"
)

func TestErrHandle(t *testing.T) {
	server := web_copy.NewHTTPServer()

	acclogMid := access_log.NewAccessLog()
	errPage := NewErrHandle()
	errPage.AddErrPage(404, []byte("404 咯"))

	server.Use(acclogMid.Build(), errPage.Build())

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
