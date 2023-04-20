package err_handle

import (
	"github.com/soluble1/mweb"
	"github.com/soluble1/mweb/middleware/access_log"
	"testing"
)

func TestErrHandle(t *testing.T) {
	server := mweb.NewHTTPServer()

	acclogMid := access_log.NewAccessLog()
	errPage := NewErrHandle()
	errPage.AddErrPage(404, []byte("404 咯"))

	server.Use(acclogMid.Build(), errPage.Build())

	server.Get("/", func(ctx *mweb.Context) {
		ctx.Resp.Write([]byte("this is /"))
	})

	server.Get("/user/name/age", func(ctx *mweb.Context) {
		ctx.Resp.Write([]byte("this is /user/name/age"))
	})

	server.Get("/user/name", func(ctx *mweb.Context) {
		ctx.Resp.Write([]byte("this is /user/name"))
	})

	server.Start(":8081")
}
