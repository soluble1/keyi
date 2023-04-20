package access_log

import (
	"github.com/soluble1/mweb"
	"testing"
)

func TestAccessLog(t *testing.T) {
	server := mweb.NewHTTPServer()

	acclogMid := NewAccessLog()

	server.Use(acclogMid.Build())

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
