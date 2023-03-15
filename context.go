package web_copy

import (
	"net/http"
)

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	// 参数路由
	paramPath map[string]string
	// 路由
	MatchedRoute string

	// 模板引擎
	templateEngine TemplateEngine

	// 缓存响应信息，最后刷新出去
	RespStatusCode int
	RespData       []byte
}

func (c *Context) Render(tplName string, data []byte) error {
	var err error
	c.RespData, err = c.templateEngine.Render(c.Req.Context(), tplName, data)
	c.RespStatusCode = 200
	if err != nil {
		c.RespStatusCode = 500
	}
	return err
}
