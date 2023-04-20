package recovery

import (
	"github.com/soluble1/mweb"
	"github.com/soluble1/mweb/middleware/access_log"
	"github.com/soluble1/mweb/middleware/err_handle"
	"log"
	"testing"
)

func TestRecovery(t *testing.T) {
	server := mweb.NewHTTPServer()

	acclogMid := access_log.NewAccessLog()

	errPage := err_handle.NewErrHandle()
	errPage.AddErrPage(404, []byte("404 咯"))

	recovery := NewRecoverHandle(520, []byte("你panic了！"), func(ctx *mweb.Context) {
		log.Println(ctx.Req.URL.Path)
	})

	server.Use(acclogMid.Build(), errPage.Build(), recovery.Build())

	server.Get("/", func(ctx *mweb.Context) {
		ctx.Resp.Write([]byte("this is /"))
	})

	server.Get("/user/name/age", func(ctx *mweb.Context) {
		ctx.Resp.Write([]byte("this is /user/name/age"))
	})

	server.Get("/user/name", func(ctx *mweb.Context) {
		ctx.Resp.Write([]byte("this is /user/name"))
	})

	server.Get("/panic/name", func(ctx *mweb.Context) {
		panic("panic test")
	})

	server.Start(":8081")
}
