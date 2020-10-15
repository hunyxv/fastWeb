package fastweb

import (
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

type SvcOption func(*svcOption)

type svcOption struct {
	// ErrorHandler for returning a response in case of an error while receiving or parsing the request.
	//
	// The following is a non-exhaustive list of errors that can be expected as argument:
	//   * io.EOF
	//   * io.ErrUnexpectedEOF
	//   * ErrGetOnly
	//   * ErrSmallBuffer
	//   * ErrBodyTooLarge
	//   * ErrBrokenChunks
	ErrorHandler func(ctx *fasthttp.RequestCtx, err error)

	// HeaderReceived is called after receiving the header
	//
	// non zero RequestConfig field values will overwrite the default configs
	HeaderReceived func(header *fasthttp.RequestHeader) fasthttp.RequestConfig

	// ContinueHandler is called after receiving the Expect 100 Continue Header
	//
	// https://www.w3.org/Protocols/rfc2616/rfc2616-sec8.html#sec8.2.3
	// https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.1.1
	// Using ContinueHandler a server can make decisioning on whether or not
	// to read a potentially large request body based on the headers
	//
	// The default is to automatically read request bodies of Expect 100 Continue requests
	// like they are normal requests
	ContinueHandler func(header *fasthttp.RequestHeader) bool

	// Server name for sending in response headers.
	//
	// Default server name is used if left blank.
	Name string

	// The maximum number of concurrent connections the server may serve.
	//
	// fasthttp.DefaultConcurrency is used if not set.
	//
	// Concurrency only works if you either call Serve once, or only ServeConn multiple times.
	// It works with ListenAndServe as well.
	Concurrency int

	// Whether to disable keep-alive connections.
	//
	// The server will close all the incoming connections after sending
	// the first response to client if this option is set to true.
	//
	// By default keep-alive connections are enabled.
	DisableKeepalive bool

	// Per-connection buffer size for requests' reading.
	// This also limits the maximum header size.
	//
	// Increase this buffer if your clients send multi-KB RequestURIs
	// and/or multi-KB headers (for example, BIG cookies).
	//
	// Default buffer size is used if not set.
	ReadBufferSize int

	// Per-connection buffer size for responses' writing.
	//
	// Default buffer size is used if not set.
	WriteBufferSize int

	// ReadTimeout is the amount of time allowed to read
	// the full request including body. The connection's read
	// deadline is reset when the connection opens, or for
	// keep-alive connections after the first byte has been read.
	//
	// By default request read timeout is unlimited.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset after the request handler
	// has returned.
	//
	// By default response write timeout is unlimited.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alive is enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used.
	IdleTimeout time.Duration

	// Maximum number of concurrent client connections allowed per IP.
	//
	// By default unlimited number of concurrent connections
	// may be established to the server from a single IP address.
	MaxConnsPerIP int

	// Maximum number of requests served per connection.
	//
	// The server closes connection after the last request.
	// 'Connection: close' header is added to the last response.
	//
	// By default unlimited number of requests may be served per connection.
	MaxRequestsPerConn int

	// MaxKeepaliveDuration is a no-op and only left here for backwards compatibility.
	// Deprecated: Use IdleTimeout instead.
	MaxKeepaliveDuration time.Duration

	// Whether to enable tcp keep-alive connections.
	//
	// Whether the operating system should send tcp keep-alive messages on the tcp connection.
	//
	// By default tcp keep-alive connections are disabled.
	TCPKeepalive bool

	// Period between tcp keep-alive messages.
	//
	// TCP keep-alive period is determined by operation system by default.
	TCPKeepalivePeriod time.Duration

	// Maximum request body size.
	//
	// The server rejects requests with bodies exceeding this limit.
	//
	// Request body size is limited by fasthttp.DefaultMaxRequestBodySize by default.
	MaxRequestBodySize int

	// Aggressively reduces memory usage at the cost of higher CPU usage
	// if set to true.
	//
	// Try enabling this option only if the server consumes too much memory
	// serving mostly idle keep-alive connections. This may reduce memory
	// usage by more than 50%.
	//
	// Aggressive memory usage reduction is disabled by default.
	ReduceMemoryUsage bool

	// Rejects all non-GET requests if set to true.
	//
	// This option is useful as anti-DoS protection for servers
	// accepting only GET requests. The request size is limited
	// by ReadBufferSize if GetOnly is set.
	//
	// Server accepts all the requests by default.
	GetOnly bool

	// Will not pre parse Multipart Form data if set to true.
	//
	// This option is useful for servers that desire to treat
	// multipart form data as a binary blob, or choose when to parse the data.
	//
	// Server pre parses multipart form data by default.
	DisablePreParseMultipartForm bool

	// Logs all errors, including the most frequent
	// 'connection reset by peer', 'broken pipe' and 'connection timeout'
	// errors. Such errors are common in production serving real-world
	// clients.
	//
	// By default the most frequent errors such as
	// 'connection reset by peer', 'broken pipe' and 'connection timeout'
	// are suppressed in order to limit output log traffic.
	LogAllErrors bool

	// Header names are passed as-is without normalization
	// if this option is set.
	//
	// Disabled header names' normalization may be useful only for proxying
	// incoming requests to other servers expecting case-sensitive
	// header names. See https://github.com/valyala/fasthttp/issues/57
	// for details.
	//
	// By default request and response header names are normalized, i.e.
	// The first letter and the first letters following dashes
	// are uppercased, while all the other letters are lowercased.
	// Examples:
	//
	//     * HOST -> Host
	//     * content-type -> Content-Type
	//     * cONTENT-lenGTH -> Content-Length
	DisableHeaderNamesNormalizing bool

	// SleepWhenConcurrencyLimitsExceeded is a duration to be slept of if
	// the concurrency limit in exceeded (default [when is 0]: don't sleep
	// and accept new connections immidiatelly).
	SleepWhenConcurrencyLimitsExceeded time.Duration

	// NoDefaultServerHeader, when set to true, causes the default Server header
	// to be excluded from the Response.
	//
	// The default Server header value is the value of the Name field or an
	// internal default value in its absence. With this option set to true,
	// the only time a Server header will be sent is if a non-zero length
	// value is explicitly provided during a request.
	NoDefaultServerHeader bool

	// NoDefaultDate, when set to true, causes the default Date
	// header to be excluded from the Response.
	//
	// The default Date header value is the current date value. When
	// set to true, the Date will not be present.
	NoDefaultDate bool

	// NoDefaultContentType, when set to true, causes the default Content-Type
	// header to be excluded from the Response.
	//
	// The default Content-Type header value is the internal default value. When
	// set to true, the Content-Type will not be present.
	NoDefaultContentType bool

	// ConnState specifies an optional callback function that is
	// called when a client connection changes state. See the
	// ConnState type and associated constants for details.
	ConnState func(net.Conn, fasthttp.ConnState) // 钩子当连接状态发生变化时执行

	// Logger, which is used by RequestCtx.Logger().
	//
	// By default standard logger from log package is used.
	Logger fasthttp.Logger

	// KeepHijackedConns is an opt-in disable of connection
	// close by fasthttp after connections' HijackHandler returns.
	// This allows to save goroutines, e.g. when fasthttp used to upgrade
	// http connections to WS and connection goes to another handler,
	// which will close it when needed.
	KeepHijackedConns bool
}

// WithErrorHandler set ErrorHandler
func WithErrorHandler(errHandler func(ctx *fasthttp.RequestCtx, err error)) SvcOption {
	return func(opt *svcOption) {
		opt.ErrorHandler = errHandler
	}
}

// WithHeaderReceived set ReadTimeout,WriteTimeout,MaxRequestBodySize options based on header
// which overrides default config
func WithHeaderReceived(headerReceived func(header *fasthttp.RequestHeader) fasthttp.RequestConfig) SvcOption {
	return func(opt *svcOption) {
		opt.HeaderReceived = headerReceived
	}
}

// WithContinueHandler set ContinueHandler
// ContinueHandler is called after receiving the Expect 100 Continue Header
// Using ContinueHandler a server can make decisioning on whether or not
// to read a potentially large request body based on the headers
func WithContinueHandler(continueHandler func(header *fasthttp.RequestHeader) bool) SvcOption {
	return func(opt *svcOption) {
		opt.ContinueHandler = continueHandler
	}
}

// WithName set server name
func WithName(name string) SvcOption {
	return func(opt *svcOption) {
		opt.Name = name
	}
}

// WithConcurrency set Concurrency
func WithConcurrency(concurrency int) SvcOption {
	return func(opt *svcOption) {
		opt.Concurrency = concurrency
	}
}

// WithDisableKeepalive disable keepalive
func WithDisableKeepalive() SvcOption {
	return func(opt *svcOption) {
		opt.DisableKeepalive = true
	}
}

// WithReadBufferSize Per-connection buffer size for requests' reading.
func WithReadBufferSize(size int) SvcOption {
	return func(opt *svcOption) {
		opt.ReadBufferSize = size
	}
}

// WithWriteBufferSize Per-connection buffer size for responses' writing.
func WithWriteBufferSize(size int) SvcOption {
	return func(opt *svcOption) {
		opt.WriteBufferSize = size
	}
}

// WithReadTimeout ReadTimeout is the amount of time
// allowed to read the full request including body.
func WithReadTimeout(readTimeout time.Duration) SvcOption {
	return func(opt *svcOption) {
		opt.ReadTimeout = readTimeout
	}
}

// WithWriteTimeout WriteTimeout is the maximum duration
// before timing out writes of the response.
func WithWriteTimeout(writeTimeout time.Duration) SvcOption {
	return func(opt *svcOption) {
		opt.WriteTimeout = writeTimeout
	}
}

// WithIdleTimeout set IdleTimeout
func WithIdleTimeout(idleTimeout time.Duration) SvcOption {
	return func(opt *svcOption) {
		opt.IdleTimeout = idleTimeout
	}
}

// WithMaxConnsPerIP set MaxConnsPerIP
func WithMaxConnsPerIP(maximum int) SvcOption {
	return func(opt *svcOption) {
		opt.MaxConnsPerIP = maximum
	}
}

// WithMaxRequestsPerConn set MaxRequestsPerConn
func WithMaxRequestsPerConn(maximum int) SvcOption {
	return func(opt *svcOption) {
		opt.MaxRequestsPerConn = maximum
	}
}

// WithMaxKeepaliveDuration set MaxKeepaliveDuration
func WithMaxKeepaliveDuration(maxkt time.Duration) SvcOption {
	return func(opt *svcOption) {
		opt.MaxKeepaliveDuration = maxkt
	}
}

// WithTCPKeepalive Whether the operating system should 
// send tcp keep-alive messages on the tcp connection.
func WithTCPKeepalive() SvcOption {
	return func(opt *svcOption) {
		opt.TCPKeepalive = true
	}
}

// WithTCPKeepalivePeriod set the time interval between tcp keep alive messages.
func WithTCPKeepalivePeriod(period time.Duration) SvcOption {
	return func(opt *svcOption) {
		opt.TCPKeepalivePeriod = period
	}
}