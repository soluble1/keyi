package recovery

import (
	"log"
	"testing"
	"web_copy"
	"web_copy/middleware/access_log"
	"web_copy/middleware/err_handle"
)

func TestRecovery(t *testing.T) {
	server := web_copy.NewHTTPServer()

	acclogMid := access_log.NewAccessLog()

	errPage := err_handle.NewErrHandle()
	errPage.AddErrPage(404, []byte("404 咯"))

	recovery := NewRecoverHandle(520, []byte("你panic了！"), func(ctx *web_copy.Context) {
		log.Println(ctx.Req.URL.Path)
	})

	server.Use(acclogMid.Build(), errPage.Build(), recovery.Build())

	server.Get("/", func(ctx *web_copy.Context) {
		ctx.Resp.Write([]byte("this is /"))
	})

	server.Get("/user/name/age", func(ctx *web_copy.Context) {
		ctx.Resp.Write([]byte("this is /user/name/age"))
	})

	server.Get("/user/name", func(ctx *web_copy.Context) {
		ctx.Resp.Write([]byte("this is /user/name"))
	})

	server.Get("/panic/name", func(ctx *web_copy.Context) {
		panic("panic test")
	})

	server.Start(":8081")
}
