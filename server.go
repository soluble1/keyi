package mweb

import (
	"log"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error
	addRouter(method string, path string, handlerFunc HandlerFunc)
}

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	router
	mdls           []Middleware
	templateEngine TemplateEngine
}

func WithGoTemplateEngine(engine TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.templateEngine = engine
	}
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *HTTPServer) Use(mdls ...Middleware) {
	if s.mdls == nil {
		s.mdls = mdls
		return
	}
	s.mdls = append(s.mdls, mdls...)
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:            request,
		Resp:           writer,
		templateEngine: s.templateEngine,
	}

	root := s.serve
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}

	var m = Middleware(func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flushResq(ctx)
		}
	})

	root = m(root)
	root(ctx)
}

func (s *HTTPServer) flushResq(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("写响应失败", err)
	}
}

func (s *HTTPServer) Get(addr string, handlerFunc HandlerFunc) {
	s.addRouter(http.MethodGet, addr, handlerFunc)
}

func (s *HTTPServer) Post(addr string, handlerFunc HandlerFunc) {
	s.addRouter(http.MethodPost, addr, handlerFunc)
}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRouter(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handlerFunc == nil {
		ctx.RespStatusCode = 404
		ctx.RespData = []byte("404 not found")
		return
	}

	ctx.MatchedRoute = mi.n.route
	ctx.paramPath = mi.paramPath
	mi.n.handlerFunc(ctx)
}
