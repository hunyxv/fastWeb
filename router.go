package fastweb

import (
	"fmt"
	"log"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(ctx *Context) {
	key := fmt.Sprintf("%s-%s", ctx.Method(), ctx.RequestURI())
	if handler, ok := r.handlers[key]; ok {
		handler(ctx)
	} else {
		fmt.Fprintf(ctx, "404 NOT FOUND: %s\n", ctx.RequestURI())
	}
}