package mweb

import "testing"

func TestRooter(t *testing.T) {
	s := NewHTTPServer()

	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("this is /"))
	})

	s.Get("/name", func(ctx *Context) {
		ctx.Resp.Write([]byte("this is /name"))
	})

	s.Get("/user/age/:id", func(ctx *Context) {
		ctx.Resp.Write([]byte("this is /user/age/:id"))
		ctx.Resp.Write([]byte(ctx.paramPath["id"]))
	})

	s.Get("/user/:id", func(ctx *Context) {
		ctx.RespData = []byte("user/:id")
	})

	s.Get("/user/:id/a/*", func(ctx *Context) {
		ctx.RespData = []byte("user/:id/a/*")
	})

	s.Start(":8081")
}
