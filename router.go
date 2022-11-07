package gweb

import (
	"log"
	"strings"
)

// node 路由节点
// 每一个节点对应路径里面的一段，组成一个树形结构
// 比如路径/a/b，a是一个node，b也是一个node
type node struct {
	//路由 path
	pattern string
	// 路由path里面的一段
	part     string
	children []*node
	// 标记节点是否是匹配参数的路径，类似:id、:name的路径节点
	isWild bool
}

// router 路由对象
// roots和handlers的key是一样的，roots匹配成功之后，在handlers里面取到对应的处理函数处理请求
type router struct {
	// 路由匹配节点,key是http method，
	roots map[string]*node
	// 路由处理函数，key是method+path
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

// AddRoute 添加路由
// method：http method  GET POST...
// pattern: 匹配路径  /a/b    /a/abc/gh
func (r *router) AddRoute(method, pattern string, handler HandleFunc) {
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{
			part: parts[0],
		}
	}
	//r.roots 的每一个元素都是一颗路由节点树
	r.roots[method].insert(pattern, parts, 1)
	r.handlers[method+"-"+pattern] = handler
}

// Match 匹配路由
// 返回节点，和url参数
func (r *router) Match(method, path string) (*node, map[string]string) {
	pathParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, params
	}

	// 匹配url path
	n := root.search(pathParts, 1)

	if n == nil {
		return nil, nil
	}
	// 解析path参数
	parts := parsePattern(n.pattern)
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = pathParts[index]
		}
		if part[0] == '*' && len(path) > 1 {
			params[part[1:]] = strings.Join(pathParts[index:], "/")
		}
	}
	return n, params
}

// parsePattern 把路径按照字符/分隔出来
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for key, item := range vs {
		if key > 0 && item == "" {
			continue
		}
		if key == 0 {
			//第一个节点是根节点
			item = "/"
		}
		parts = append(parts, item)
	}
	return parts
}

// insert 插入路由节点
// 顶部节点是根，先从第一个节点开始
// height：节点高度
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	childNode := n.findChild(part)
	if childNode == nil {
		childNode = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, childNode)
	}
	childNode.insert(pattern, parts, height+1)
}

// search 搜索节点
// 根据url path的节点一直往下搜索到最后一个,全部匹配说明这个路由匹配成功
func (n *node) search(parts []string, height int) *node {
	// 根据path的节点，已经匹配到最后一个节点了，pattern不为空则匹配成功，否则失败
	// 上面insert的时候，路由path的最后一个节点才会设置pattern的值，之前的节点是为空的
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//只要判断pattern的值是否为空就知道搜索的链条是否和url path匹配
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.findAllChild(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// findChild insert的时候查找节点的子节点
func (n *node) findChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// findAllChild 查找所有匹配的子节点
func (n *node) findAllChild(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 路由分组
type RouterGroup struct {
	prefix string
	// 中间件挂在路由分组上
	middlewares []HandleFunc
	parent      *RouterGroup
	engine      *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
	}
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

// addRoute 添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	// 分组添加路由，路由匹配模式是分组路径加上当前路径
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.AddRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
