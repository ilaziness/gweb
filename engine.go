package gweb

import (
	"log"
	"net/http"
)

type HandleFunc func(*Context)

func defaultHandler(g *Context) {
	g.String(http.StatusOK, "web service is running.")
}

// Engine 框架核心
type Engine struct {
	router *router
}

func NewDefault() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

// Run 启动运行服务
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP 实现Handler接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.handle(newContext(w, req))
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
}

// handle 匹配路由，执行handler
func (e *Engine) handle(g *Context) {
	log.Println(g.Path)
	n, params := e.router.Match(g.Method, g.Path)
	if n == nil && g.Path == "/" {
		defaultHandler(g)
		return
	}
	if n == nil {
		g.String(http.StatusNotFound, "404 NOT FOUND")
		return
	}
	g.Params = params
	key := g.Method + "-" + n.pattern
	e.router.handlers[key](g)
}

// AddRoute 添加路由记录
func (e *Engine) AddRoute(method, pattern string, handler HandleFunc) {
	e.router.AddRoute(method, pattern, handler)
}

// GET 添加GET请求路由
func (e *Engine) GET(pattern string, handler HandleFunc) {
	e.router.AddRoute(http.MethodGet, pattern, handler)
}
