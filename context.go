package web_copy

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
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

	// 缓存的请求中的数据 /user?name=name&age=1
	cacheQueryValue url.Values

	// 缓存响应信息，最后刷新出去
	RespStatusCode int
	RespData       []byte
}

// Render 渲染模板
func (c *Context) Render(tplName string, data []byte) error {
	var err error
	c.RespData, err = c.templateEngine.Render(c.Req.Context(), tplName, data)
	c.RespStatusCode = 200
	if err != nil {
		c.RespStatusCode = 500
	}
	return err
}

// BindJSON 解析body
func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("body 为空")
	}
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.RespStatusCode = code
	c.RespData = bs
	return err
}

// FormValue 获取表单中 key 对应的 value 值
func (c *Context) FormValue(key string) StringValue {
	// ParseForm 内部判断了是否重复解析问题
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{val: "", err: err}
	}
	return StringValue{val: c.Req.FormValue(key), err: nil}
}

// QueryValue 查询路径中的参数
func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValue == nil {
		c.cacheQueryValue = c.Req.URL.Query()
	}
	val, ok := c.cacheQueryValue[key]
	if !ok && len(c.cacheQueryValue[key]) == 0 {
		return StringValue{val: "", err: errors.New("web：找不到这个key")}
	}
	return StringValue{val: val[0], err: nil}
}

// PathValue 查询参数路径
func (c *Context) PathValue(key string) StringValue {
	val, ok := c.paramPath[key]
	if !ok {
		return StringValue{val: "", err: errors.New("web：找不到这个key")}
	}
	return StringValue{val: val, err: nil}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

func (s StringValue) ToInt32() (int32, error) {
	if s.err != nil {
		return 0, s.err
	}
	val, err := strconv.ParseInt(s.val, 10, 32)

	return int32(val), err
}
