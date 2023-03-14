package web_copy

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error
	addRouter(method string, path string, handlerFunc HandlerFunc)
}

type HTTPServer struct {
	router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	s.serve(ctx)
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
		return
	}

	ctx.paramPath = mi.paramPath
	mi.n.handlerFunc(ctx)
}
