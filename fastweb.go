package fastweb

import (
	// "reflect"
	"github.com/valyala/fasthttp"
)

type HandlerFunc func(ctx *Context)

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

func (engine *Engine) requestHandler (fctx *fasthttp.RequestCtx) {
	ctx := NewContext(fctx)
	engine.router.handle(ctx)
	ReleaseCtx(ctx)
}