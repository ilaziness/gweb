package gweb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	MimeTypeJson = "application/json;charset=utf-8"
	MimeTypeHtml = "text/html;charset=utf-8"
)

type Context struct {
	ctx        context.Context
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	// url匹配到的参数，包含路由参数和url query的参数
	Params map[string]string
	//middleware,中间件也是一个普通的handler
	handlers []HandleFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// SetHeader 设置http头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// Status 设置http status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// String 返回文本响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain;charset=utf-8")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 响应json数据
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", MimeTypeJson)
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		c.RespServerError(err.Error())
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// RespServerError 返回500错误
func (c *Context) RespServerError(err string) {
	http.Error(c.Writer, err, http.StatusInternalServerError)
}

// Next 执行下一个中间件
// 依次执行c.handlers待执行的中间件和路由handler
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
