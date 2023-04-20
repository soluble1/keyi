package mweb

import "testing"

func TestTemplateEngine(t *testing.T) {
	g1 := &GoTemplateEngine{}

	server := NewHTTPServer(WithGoTemplateEngine(g1))
	server.Get("/test", func(ctx *Context) {
		g1.LoadFromGlob("./testData/tpls/test.html")

		er := ctx.Render("test.html", nil)
		if er != nil {
			t.Fatal(er)
		}
	})

	server.Get("/login", func(ctx *Context) {
		g1.LoadFromGlob("./testData/tpls/login.gohtml")

		er := ctx.Render("login.gohtml", nil)
		if er != nil {
			t.Fatal(er)
		}
	})

	server.Get("/test2", func(ctx *Context) {
		g1.LoadFromGlob("./testData/tpls/test2.html")

		er := ctx.Render("test2.html", nil)
		if er != nil {
			t.Fatal(er)
		}
	})

	err := server.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}
