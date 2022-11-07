package gweb

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(*Context)

func defaultHandler(g *Context) {
	g.String(http.StatusOK, "web service is running.")
}

// Engine 框架核心
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func NewDefault() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Run 启动运行服务
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP 实现Handler接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleware []HandleFunc
	// 找到所有需要应用的中间件
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleware = append(middleware, group.middlewares...)
		}
	}
	c := newContext(w, req)
	// 需要应用的中间件放到请求的handlers里面待执行
	c.handlers = middleware
	e.handle(c)
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
}

// handle 匹配路由，执行handler
func (e *Engine) handle(g *Context) {
	log.Println(g.Path)
	n, params := e.router.Match(g.Method, g.Path)
	if n == nil && g.Path == "/" {
		// 默认路由
		g.handlers = append(g.handlers, defaultHandler)
	} else if n == nil {
		// 404
		g.handlers = append(g.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND")
		})
	} else {
		g.Params = params
		key := g.Method + "-" + n.pattern
		// 路由handler，放到待执行列表最后
		g.handlers = append(g.handlers, e.router.handlers[key])
	}
	// 开始执行中间件
	g.Next()
}

// AddRoute 添加路由记录
func (e *Engine) AddRoute(method, pattern string, handler HandleFunc) {
	e.router.AddRoute(method, pattern, handler)
}

// GET 添加GET请求路由
func (e *Engine) GET(pattern string, handler HandleFunc) {
	e.router.AddRoute(http.MethodGet, pattern, handler)
}
