package test

import (
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	mweb "github.com/soluble1/mweb"
	"github.com/soluble1/mweb/session"
	"github.com/soluble1/mweb/session/cookie"
	mcache "github.com/soluble1/mweb/session/redis_cache"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestRedisCacheSession(t *testing.T) {
	server := mweb.NewHTTPServer()

	client := redis.NewClient(&redis.Options{
		Addr:     "120.46.196.48:6379",
		Password: "",
		DB:       0,
	})
	s := mcache.NewRedisStore(30*time.Minute, client)

	manager := session.Manager{
		SessCtxKey: "my_session",
		Store:      s,
		Propagator: cookie.NewPropagator("sessId", func(p *cookie.Propagator) {
			cookie.WithCookieOpt(func(cookie *http.Cookie) {
				// cookie.HttpOnly = true, js脚本将无法读取到cookie信息
				cookie.HttpOnly = true
			})
		}),
	}

	server.Post("/login", func(ctx *mweb.Context) {
		ctx.RespData = []byte("this is /login")

		id := uuid.New().String()
		sess, err := manager.InitSession(ctx, id)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		err = sess.Set(ctx.Req.Context(), "xiao", "xiao long ren")
		err = sess.Set(ctx.Req.Context(), "ma", "ma jun da shi")
		err = sess.Set(ctx.Req.Context(), "key", id)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
	})

	server.Get("/resource", func(ctx *mweb.Context) {
		sess, err := manager.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}

		val, err := sess.Get(ctx.Req.Context(), "key")
		if err != nil {
			log.Println(err)
		}
		ctx.RespData = []byte(val)
		ctx.RespStatusCode = 200
	})

	server.Post("/logout", func(ctx *mweb.Context) {
		_ = manager.RemoveSession(ctx)
	})

	server.Use(func(next mweb.HandlerFunc) mweb.HandlerFunc {
		return func(ctx *mweb.Context) {
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
