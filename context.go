package fastweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)


type H map[string]interface{}

type Context struct {
	fastctx		*fasthttp.RequestCtx
	Path		string
	UrlParams	map[string]string
}

var ctxPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

func NewContext(fctx *fasthttp.RequestCtx) *Context {
	ctx := ctxPool.Get().(*Context)
	ctx.fastctx = fctx 
	ctx.Path = string(fctx.RequestURI())
	return ctx
}

func ReleaseCtx(ctx *Context){
	ctx.Path = ""
	ctx.fastctx = nil
	ctxPool.Put(ctx)
}

func (c *Context) UrlParam(key string) string {
	value, _ := c.UrlParams[key]
	return value
}

func (ctx *Context) PostForm(key string) string {
	return string(ctx.fastctx.PostArgs().Peek(key))
}

func (ctx *Context) Query(key string) string {
	return string(ctx.fastctx.QueryArgs().Peek(key))
}

func (ctx *Context) Status(code int) {
	ctx.fastctx.Response.SetStatusCode(code)
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.fastctx.Response.Header.Set(key, value)
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	fmt.Fprintf(ctx.fastctx, fmt.Sprintf(format, values...))
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.fastctx)
	if err := encoder.Encode(obj); err != nil {
		// http.Error(ctx, err.Error(), 500)
		fmt.Fprintf(ctx.fastctx, err.Error())
		ctx.Status(500)
	}
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	dReader := bytes.NewReader(data)
	dReader.WriteTo(ctx.fastctx)
}

func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	fmt.Fprint(ctx.fastctx, html)
}