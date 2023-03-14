package web_copy

import "testing"

func TestRooter(t *testing.T) {
	s := NewHTTPServer()

	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("this is /"))
	})

	s.Get("/name", func(ctx *Context) {
		ctx.Resp.Write([]byte("this is /name"))
	})

	s.Get("/user/*/name", func(ctx *Context) {
		ctx.Resp.Write([]byte("this is /user/*/name"))
	})

	s.Get("/user/age/:id", func(ctx *Context) {
		ctx.Resp.Write([]byte("this is /user/age/:id"))
		ctx.Resp.Write([]byte(ctx.paramPath["id"]))
	})

	s.Start(":8081")
}
