package fastweb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)

type H map[string]interface{}

type Context interface {
	SetUserValue(key string, value interface{})
	SetUserValueBytes(key []byte, value interface{})
	UserValue(key string) interface{}
	UserValueBytes(key []byte) interface{}
	VisitUserValues(visitor func([]byte, interface{}))
	// IsTLS() bool
	// TLSConnectionState() *tls.ConnectionState
	// Conn() net.Conn
	// String() string
	// ID() uint64
	// ConnID() uint64
	// Time() time.Time
	// ConnRequestNum() uint64
	// ConnTime() time.Time
	// SetConnectionClose()
	SetStatusCode(statusCode int)
	SetContentType(contentType string)
	SetContentTypeBytes(contentType []byte)
	RequestURI() []byte
	URI() *fasthttp.URI
	Referer() []byte
	UserAgent() []byte
	Path() string
	Host() string
	// QueryArgs() *fasthttp.Args
	// PostArgs() *fasthttp.Args
	// MultipartForm() (*multipart.Form, error)
	// FormFile(key string) (*multipart.FileHeader, error)
	// FormValue(key string) []byte
	URLParam(key string) (string, bool)
	URLParams() map[string]string
	FormValue(key string) (string, bool)
	FormValues() map[string]string
	QueryParam(key string) (string, bool)
	QueryParams() map[string]string
	// IsGet() bool
	// IsPost() bool
	// IsPut() bool
	// IsDelete() bool
	// IsConnect() bool
	// IsOptions() bool
	// IsTrace() bool
	// IsPatch() bool
	// IsHead() bool
	Method() string

	// LocalAddr() net.Addr
	// RemoteIP() net.IP
	// LocalIP() net.IP
	// Error(msg string, statusCode int)
	// Success(contentType string, body []byte)
	// SuccessString(contentType, body string)
	// Redirect(uri string, statusCode int)	// 重定向 到本服务某uri
	// RedirectBytes(uri []byte, statusCode int)
	// SetBody(body []byte)
	// SetBodyString(body string)
	ResetBody()
	SendFile(path string)
	// SendFileBytes(path []byte)
	NotFound() // 404 Page not found
	// Write(p []byte) (int, error)
	// WriteString(s string) (int, error)
	// PostBody() []byte
	// SetBodyStream(bodyStream io.Reader, bodySize int)
	// SetBodyStreamWriter(sw fasthttp.StreamWriter)
	// IsBodyStream() bool
	Logger() fasthttp.Logger
	// TimeoutErrorWithCode(msg string, statusCode int)
	// TimeoutErrorWithResponse(resp *fasthttp.Response)
	// TimeoutErrorWithResponse (resp *fasthttp.Response)
	// Init(req *fasthttp.Request, remoteAddr net.Addr, logger fasthttp.Logger) // 测试时使用
	// Value(key interface{}) interface{}  // 等同于 UserValue(key)

	SetBodyStrf(int, string, ...interface{})
	JSON(int, interface{})
}

var _ Context = (*context)(nil)

type context struct {
	Context   *fasthttp.RequestCtx
	urlParams map[string]string
}

var ctxPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func (c *context) Init(ctx *fasthttp.RequestCtx) {
	c.Context = ctx
}

// func NewContext(fctx *Context) *context {
// 	ctx := ctxPool.Get().(*Context)
// 	ctx.Content = fctx
// 	ctx.Path = string(fctx.Path())
// 	return ctx
// }

func (c *context) releaseCtx() {
	c.urlParams = nil
	c.Context = nil
	ctxPool.Put(c)
}

func (c *context) SetUserValue(key string, value interface{}) {
	c.Context.SetUserValue(key, value)
}

func (c *context) SetUserValueBytes(key []byte, value interface{}) {
	c.Context.SetUserValueBytes(key, value)
}

func (c *context) UserValue(key string) interface{} {
	return c.Context.UserValue(key)
}

func (c *context) UserValueBytes(key []byte) interface{} {
	return c.Context.UserValueBytes(key)
}

func (c *context) VisitUserValues(visitor func([]byte, interface{})) {
	c.Context.VisitUserValues(visitor)
}

func (c *context) SetStatusCode(statusCode int) {
	c.Context.SetStatusCode(statusCode)
}

func (c *context) SetContentType(contentType string) {
	c.Context.SetContentType(contentType)
}

func (c *context) SetContentTypeBytes(contentType []byte) {
	c.Context.SetContentTypeBytes(contentType)
}

func (c *context) RequestURI() []byte {
	return c.Context.RequestURI()
}

func (c *context) URI() *fasthttp.URI {
	return c.Context.URI()
}

func (c *context) Referer() []byte {
	return c.Context.Referer()
}

func (c *context) UserAgent() []byte {
	return c.Context.UserAgent()
}

func (c *context) Path() string {
	return b2s(c.Context.Path())
}

func (c *context) Host() string {
	return b2s(c.Context.Host())
}

func (c *context) Method() string {
	return b2s(c.Context.Method())
}

func (c *context) ResetBody() {
	c.Context.ResetBody()
}

func (c *context) SendFile(path string) {
	c.Context.SendFile(path)
}

func (c *context) NotFound() {
	c.Context.NotFound()
}

func (c *context) Logger() fasthttp.Logger {
	return c.Context.Logger()
}

func (c *context) URLParam(key string) (string, bool) {
	val, ok := c.urlParams[key]
	return val, ok
}
func (c *context) URLParams() map[string]string {
	return c.urlParams
}

func (c *context) FormValue(key string) (string, bool) {
	value := c.Context.PostArgs().Peek(key)
	if len(value) > 0 {
		return b2s(value), true
	}
	return "", false
}

func (c *context) FormValues() map[string]string {
	args := c.Context.PostArgs()
	form := make(map[string]string, args.Len())
	args.VisitAll(func(key, val []byte) {
		form[b2s(key)] = b2s(val)
	})
	return form
}

func (c *context) QueryParam(key string) (string, bool) {
	value := c.Context.QueryArgs().Peek(key)
	if len(value) > 0 {
		return b2s(value), true
	}
	return "", false
}

func (c *context) QueryParams() map[string]string {
	args := c.Context.QueryArgs()
	params := make(map[string]string, args.Len())
	args.VisitAll(func(key, val []byte) {
		params[b2s(key)] = b2s(val)
	})
	return params
}

func (c *context) SetStatus(code int) {
	c.Context.Response.SetStatusCode(code)
}

func (c *context) SetHeader(key, value string) {
	c.Context.Response.Header.Set(key, value)
}

func (c *context) SetBodyStrf(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	fmt.Fprintf(c.Context, fmt.Sprintf(format, values...))
}

func (c *context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Context)
	if err := encoder.Encode(obj); err != nil {
		fmt.Fprintf(c.Context, err.Error())
		c.SetStatus(500)
	}
}

func (c *context) Data(code int, data []byte) {
	c.SetStatus(code)
	dReader := bytes.NewReader(data)
	dReader.WriteTo(c.Context)
}

func (c *context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	fmt.Fprint(c.Context, html)
}
