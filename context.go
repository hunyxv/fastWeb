package fastweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
)


type H map[string]interface{}
type Context struct{
	fasthttp.RequestCtx
}

func (ctx *Context) PostForm(key string) string {
	return string(ctx.PostArgs().Peek(key))
}

func (ctx *Context) Query(key string) string {
	return string(ctx.QueryArgs().Peek(key))
}

func (ctx *Context) Status(code int) {
	ctx.Response.SetStatusCode(code)
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Response.Header.Set(key, value)
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	fmt.Fprintf(ctx, fmt.Sprintf(format, values...))
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx, err.Error(), 500)
	}
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	dReader := bytes.NewReader(data)
	ctx.Response.Read(dReader)
}

func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	fmt.Fprint(ctx, html)
}