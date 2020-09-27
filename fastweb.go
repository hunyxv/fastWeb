package fastweb

import (
	"github.com/valyala/fasthttp"
)

// HandlerFunc 定义 http 处理器
type HandlerFunc func(ctx Context)

// Engine fastweb 引擎
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
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
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET 添加 Get 路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 添加 POST 路由
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
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

// GET 添加 Get 路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 添加 POST 路由
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) requestHandler(fctx *fasthttp.RequestCtx) {
	ctx := ctxPool.Get().(*context)
	ctx.Init(fctx)
	engine.router.handle(ctx)
	ctx.releaseCtx()
}

// Run 启动服务
func (engine *Engine) Run(addr string) (err error) {
	return fasthttp.ListenAndServe(addr, engine.requestHandler)
}
