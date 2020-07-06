package fastweb

import (
	"github.com/valyala/fasthttp"
)

type HandlerFunc func(ctx *fasthttp.RequestCtx)

type Engine struct {
	router	*router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET (pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST (pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run (addr string) (err error) {
	return fasthttp.ListenAndServe(addr, engine.requestHandler)
}

func (engine *Engine) requestHandler (ctx *fasthttp.RequestCtx) {
	engine.router.handle(ctx)
}