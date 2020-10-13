package fastweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)

// type H map[string]interface{}

// Context fastweb 应用上下文
type Context interface {
	GetFctx() *fasthttp.RequestCtx
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
	SetPath(string)
	Host() string
	// QueryArgs() *fasthttp.Args
	// PostArgs() *fasthttp.Args
	// MultipartForm() (*multipart.Form, error)
	// FormFile(key string) (*multipart.FileHeader, error)
	// FormValue(key string) []byte
	SetURLParam(Params)
	URLParam(string) (string, bool)
	URLParams() map[string]string
	FormValue(key string) (string, bool)
	FormValues() map[string]string
	QueryParam(string) (string, bool)
	QueryParams(interface{}) error
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
	Error(msg string, statusCode int)
	// Success(contentType string, body []byte)
	// SuccessString(contentType, body string)
	Redirect(uri string, statusCode int)	// 重定向 到本服务某uri
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
	fctx      *fasthttp.RequestCtx
	urlParams map[string]string
}

var ctxPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func (c *context) Init(ctx *fasthttp.RequestCtx) {
	c.fctx = ctx
	c.urlParams = make(map[string]string)
}

func (c *context) GetFctx() *fasthttp.RequestCtx {
	return c.fctx
}

// func NewContext(fctx *Context) *context {
// 	ctx := ctxPool.Get().(*Context)
// 	ctx.Content = fctx
// 	ctx.Path = string(fctx.Path())
// 	return ctx
// }

func (c *context) releaseCtx() {
	c.urlParams = nil
	c.fctx = nil
	ctxPool.Put(c)
}

func (c *context) SetURLParam(ps Params) {
	for _, param := range ps {
		c.urlParams[param.Key] = param.Value
	}
}

func (c *context) SetUserValue(key string, value interface{}) {
	c.fctx.SetUserValue(key, value)
}

func (c *context) SetUserValueBytes(key []byte, value interface{}) {
	c.fctx.SetUserValueBytes(key, value)
}

func (c *context) UserValue(key string) interface{} {
	return c.fctx.UserValue(key)
}

func (c *context) UserValueBytes(key []byte) interface{} {
	return c.fctx.UserValueBytes(key)
}

func (c *context) VisitUserValues(visitor func([]byte, interface{})) {
	c.fctx.VisitUserValues(visitor)
}

func (c *context) SetStatusCode(statusCode int) {
	c.fctx.SetStatusCode(statusCode)
}

func (c *context) SetContentType(contentType string) {
	c.fctx.SetContentType(contentType)
}

func (c *context) SetContentTypeBytes(contentType []byte) {
	c.fctx.SetContentTypeBytes(contentType)
}

func (c *context) RequestURI() []byte {
	return c.fctx.RequestURI()
}

func (c *context) URI() *fasthttp.URI {
	return c.fctx.URI()
}

func (c *context) Referer() []byte {
	return c.fctx.Referer()
}

func (c *context) UserAgent() []byte {
	return c.fctx.UserAgent()
}

func (c *context) Path() string {
	return b2s(c.fctx.Path())
}

func (c *context) SetPath(path string) {
	c.fctx.URI().SetPath(path)
}

func (c *context) Host() string {
	return b2s(c.fctx.Host())
}

func (c *context) Method() string {
	return b2s(c.fctx.Method())
}

func (c *context) Error(msg string, statusCode int) {
	c.fctx.Error(msg, statusCode)
}

func (c *context) Redirect(uri string, statusCode int) {
	c.fctx.Redirect(uri, statusCode)
}

func (c *context) ResetBody() {
	c.fctx.ResetBody()
}

func (c *context) SendFile(path string) {
	c.fctx.SendFile(path)
}

func (c *context) NotFound() {
	c.fctx.NotFound()
}

func (c *context) Logger() fasthttp.Logger {
	return c.fctx.Logger()
}

func (c *context) URLParam(key string) (string, bool) {
	val, ok := c.urlParams[key]
	return val, ok
}
func (c *context) URLParams() map[string]string {
	return c.urlParams
}

func (c *context) FormValue(key string) (string, bool) {
	value := c.fctx.PostArgs().Peek(key)
	if len(value) > 0 {
		return b2s(value), true
	}
	return "", false
}

func (c *context) FormValues() map[string]string {
	args := c.fctx.PostArgs()
	form := make(map[string]string, args.Len())
	args.VisitAll(func(key, val []byte) {
		form[b2s(key)] = b2s(val)
	})
	return form
}

func (c *context) QueryParam(key string) (string, bool) {
	value := c.fctx.QueryArgs().Peek(key)
	if len(value) > 0 {
		return b2s(value), true
	}
	return "", false
}

func (c *context) QueryParams(obj interface{}) error {
	ps, err := scan(obj)
	if err != nil {
		return err
	}

	args := c.fctx.QueryArgs()
	args.VisitAll(func(key, val []byte) {
		err = ps.padding(key, val, obj)
	})
	if err == nil {
		err = ps.valid(obj)
	}
	return err
}

func (c *context) SetStatus(code int) {
	c.fctx.Response.SetStatusCode(code)
}

func (c *context) SetHeader(key, value string) {
	c.fctx.Response.Header.Set(key, value)
}

func (c *context) SetBodyStrf(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	fmt.Fprintf(c.fctx, fmt.Sprintf(format, values...))
}

func (c *context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.fctx)
	if err := encoder.Encode(obj); err != nil {
		fmt.Fprintf(c.fctx, err.Error())
		c.SetStatus(500)
	}
}

func (c *context) Data(code int, data []byte) {
	c.SetStatus(code)
	dReader := bytes.NewReader(data)
	dReader.WriteTo(c.fctx)
}

func (c *context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	fmt.Fprint(c.fctx, html)
}
