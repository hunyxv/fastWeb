package fastweb

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type H map[string]interface{}

type Context interface {
	SetUserValue(key string, value interface{})
	SetUserValueBytes(key []byte, value interface{})
	UserValue(key string) interface{}
	UserValueBytes(key []byte) interface{}
	VisitUserValues(visitor func([]byte, interface{}))
	IsTLS() bool
	TLSConnectionState() *tls.ConnectionState
	Conn() net.Conn
	String() string
	ID() uint64
	ConnID() uint64
	Time() time.Time
	ConnRequestNum() uint64
	ConnTime() time.Time
	SetConnectionClose()
	SetStatusCode(statusCode int)
	SetContentType(contentType string)
	SetContentTypeBytes(contentType []byte)
	RequestURI() []byte
	URI() *fasthttp.URI
	Referer() []byte
	UserAgent() []byte
	Path() []byte
	Host() []byte
	QueryArgs() *fasthttp.Args
	PostArgs() *fasthttp.Args
	MultipartForm() (*multipart.Form, error)
	FormFile(key string) (*multipart.FileHeader, error)
	FormValue(key string) []byte
	IsGet() bool
	IsPost() bool
	IsPut() bool
	IsDelete() bool
	// IsConnect() bool
	IsOptions() bool
	// IsTrace() bool
	IsPatch() bool
	Method() []byte
	IsHead() bool
	// LocalAddr() net.Addr
	// RemoteIP() net.IP
	// LocalIP() net.IP
	Error(msg string, statusCode int) 
	Success(contentType string, body []byte)
	SuccessString(contentType, body string)
	Redirect(uri string, statusCode int)	// 重定向 到本服务某uri
	// RedirectBytes(uri []byte, statusCode int) 
	SetBody(body []byte)
	SetBodyString(body string)
	ResetBody() 
	SendFile(path string)
	// SendFileBytes(path []byte)
	NotFound()  // 404 Page not found
	Write(p []byte) (int, error) 
	WriteString(s string) (int, error)
	PostBody() []byte
	// SetBodyStream(bodyStream io.Reader, bodySize int)
	// SetBodyStreamWriter(sw fasthttp.StreamWriter) 
	// IsBodyStream() bool
	Logger() fasthttp.Logger
	TimeoutErrorWithCode(msg string, statusCode int)
	TimeoutErrorWithResponse(resp *fasthttp.Response)
	// TimeoutErrorWithResponse (resp *fasthttp.Response)
	// Init(req *fasthttp.Request, remoteAddr net.Addr, logger fasthttp.Logger) // 测试时使用
	// Value(key interface{}) interface{}  // 等同于 UserValue(key)
}

type context struct {
	Context *fasthttp.RequestCtx
}

var ctxPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func (c *context) Init(ctx *fasthttp.RequestCtx) {
	c.Context = ctx
}
// func NewContext(fctx *fasthttp.RequestCtx) *context {
// 	ctx := ctxPool.Get().(*Context)
// 	ctx.Content = fctx
// 	ctx.Path = string(fctx.Path())
// 	return ctx
// }

func (c *context) ReleaseCtx() {
	
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
