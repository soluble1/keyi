package web_copy

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	// 参数路由
	paramPath map[string]string
	// 路由
	MatchedRoute string

	// 缓存响应信息，最后刷新出去
	RespStatusCode int
	RespData       []byte
}
