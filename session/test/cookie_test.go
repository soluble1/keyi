package test

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"testing"
	"time"
	"web_copy"
	"web_copy/session"
	"web_copy/session/cookie"
	"web_copy/session/memory"
)

func TestCookie(t *testing.T) {
	server := web_copy.NewHTTPServer()

	manager := session.Manager{
		SessCtxKey: "my_session",
		Store:      memory.NewStore(30 * time.Minute),
		Propagator: cookie.NewPropagator("sessId", func(p *cookie.Propagator) {
			cookie.WithCookieOpt(func(cookie *http.Cookie) {
				// cookie.HttpOnly = true, js脚本将无法读取到cookie信息
				cookie.HttpOnly = true
			})
		}),
	}

	server.Post("/login", func(ctx *web_copy.Context) {
		ctx.RespData = []byte("this is /login")

		id := uuid.New()
		sess, err := manager.InitSession(ctx, id.String())
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}

		err = sess.Set(ctx.Req.Context(), "myKey", "xiaoLongRen")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
	})

	server.Get("/resource", func(ctx *web_copy.Context) {
		sess, err := manager.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}

		val, err := sess.Get(ctx.Req.Context(), "myKey")
		ctx.RespData = []byte(val)
		ctx.RespStatusCode = 200
	})

	server.Post("/logout", func(ctx *web_copy.Context) {
		_ = manager.RemoveSession(ctx)
	})

	server.Use(func(next web_copy.HandlerFunc) web_copy.HandlerFunc {
		return func(ctx *web_copy.Context) {
			if ctx.Req.URL.Path != "/login" {
				sess, err := manager.GetSession(ctx)
				if err != nil {
					ctx.RespStatusCode = 401
					log.Println(err)
					return
				}
				_ = manager.Refresh(ctx.Req.Context(), sess.ID())
			}
			next(ctx)
		}
	})

	server.Start(":8081")
}
