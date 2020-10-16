package fastweb

import (
	"github.com/valyala/fasthttp"
)

// HandlerFunc 定义 http 处理器
type HandlerFunc func(ctx Context)

// Engine fastweb 引擎
type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
	logger fasthttp.Logger
}

// RouterGroup 路由分组结构体
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

// New 返回 *Engine 实例
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// addRoute 向根 group 添加路由
func (engine *Engine) addRoute(method string, pattern string, handle HandlerFunc) {
	engine.router.addRoute(method, pattern, handle)
}

func (engine *Engine) GET(pattern string, handle HandlerFunc) {
	engine.addRoute(fasthttp.MethodGet, pattern, handle)
}

func (engine *Engine) HEAD(path string, handle HandlerFunc) {
	engine.addRoute(fasthttp.MethodHead, path, handle)
}

func (engine *Engine) OPTIONS(path string, handle HandlerFunc) {
	engine.addRoute(fasthttp.MethodOptions, path, handle)
}

func (engine *Engine) POST(pattern string, handle HandlerFunc) {
	engine.addRoute(fasthttp.MethodPost, pattern, handle)
}

func (engine *Engine) PUT(path string, handle HandlerFunc) {
	engine.addRoute(fasthttp.MethodPut, path, handle)
}

func (engine *Engine) PATCH(path string, handle HandlerFunc) {
	engine.addRoute(fasthttp.MethodPatch, path, handle)
}

func (engine *Engine) DELETE(path string, handle HandlerFunc) {
	engine.addRoute(fasthttp.MethodDelete, path, handle)
}

func (engine *Engine) ServeFiles(path, root string) {
	engine.router.ServeFiles(path, root)
}

// Group 创建分组路由
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 向当前分组路由中添加路由
func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	// log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handle HandlerFunc) {
	group.addRoute(fasthttp.MethodGet, pattern, handle)
}

func (group *RouterGroup) HEAD(path string, handle HandlerFunc) {
	group.addRoute(fasthttp.MethodHead, path, handle)
}

func (group *RouterGroup) OPTIONS(path string, handle HandlerFunc) {
	group.addRoute(fasthttp.MethodOptions, path, handle)
}

func (group *RouterGroup) POST(pattern string, handle HandlerFunc) {
	group.addRoute(fasthttp.MethodPost, pattern, handle)
}

func (group *RouterGroup) PUT(path string, handle HandlerFunc) {
	group.addRoute(fasthttp.MethodPut, path, handle)
}

func (group *RouterGroup) PATCH(path string, handle HandlerFunc) {
	group.addRoute(fasthttp.MethodPatch, path, handle)
}

func (group *RouterGroup) DELETE(path string, handle HandlerFunc) {
	group.addRoute(fasthttp.MethodDelete, path, handle)
}

func (group *RouterGroup) ServeFiles(path, root string) {
	group.engine.router.ServeFiles(path, root)
}

func (engine *Engine) requestHandler(fctx *fasthttp.RequestCtx) {
	ctx := ctxPool.Get().(*context)
	ctx.Init(fctx)
	engine.router.Handle(ctx)
	ctx.releaseCtx()
}

// Run 启动服务
func (engine *Engine) Run(addr string, options ...SvrOption) (err error) {
	server := &fasthttp.Server{
		Handler: engine.requestHandler,
		Name:    "fastweb",
		Logger:  engine.logger,
	}

	for _, f := range options {
		f(server)
	}

	return server.ListenAndServe(addr) // fasthttp.ListenAndServe(addr, engine.requestHandler)
}
