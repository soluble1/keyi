package web_copy

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	// 参数路由
	paramPath map[string]string
}
