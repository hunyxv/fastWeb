package fastweb

import (
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

// SvrOption set fasthttp.Server option
type SvrOption func(*fasthttp.Server)

// WithErrorHandler set ErrorHandler
func WithErrorHandler(errHandler func(ctx *fasthttp.RequestCtx, err error)) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.ErrorHandler = errHandler
	}
}

// WithHeaderReceived set ReadTimeout,WriteTimeout,MaxRequestBodySize options based on header
// which overrides default config
func WithHeaderReceived(headerReceived func(header *fasthttp.RequestHeader) fasthttp.RequestConfig) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.HeaderReceived = headerReceived
	}
}

// WithContinueHandler set ContinueHandler
// ContinueHandler is called after receiving the Expect 100 Continue Header
// Using ContinueHandler a server can make decisioning on whether or not
// to read a potentially large request body based on the headers
func WithContinueHandler(continueHandler func(header *fasthttp.RequestHeader) bool) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.ContinueHandler = continueHandler
	}
}

// WithName set server name
func WithName(name string) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.Name = name
	}
}

// WithConcurrency set Concurrency
func WithConcurrency(concurrency int) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.Concurrency = concurrency
	}
}

// WithDisableKeepalive disable keepalive
func WithDisableKeepalive() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.DisableKeepalive = true
	}
}

// WithReadBufferSize Per-connection buffer size for requests' reading.
func WithReadBufferSize(size int) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.ReadBufferSize = size
	}
}

// WithWriteBufferSize Per-connection buffer size for responses' writing.
func WithWriteBufferSize(size int) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.WriteBufferSize = size
	}
}

// WithReadTimeout ReadTimeout is the amount of time
// allowed to read the full request including body.
func WithReadTimeout(t time.Duration) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.ReadTimeout = t
	}
}

// WithWriteTimeout WriteTimeout is the maximum duration
// before timing out writes of the response.
func WithWriteTimeout(t time.Duration) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.WriteTimeout = t
	}
}

// WithIdleTimeout set IdleTimeout
func WithIdleTimeout(t time.Duration) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.IdleTimeout = t
	}
}

// WithMaxConnsPerIP set MaxConnsPerIP
func WithMaxConnsPerIP(maximum int) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.MaxConnsPerIP = maximum
	}
}

// WithMaxRequestsPerConn set MaxRequestsPerConn
func WithMaxRequestsPerConn(maximum int) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.MaxRequestsPerConn = maximum
	}
}

// WithMaxKeepaliveDuration set MaxKeepaliveDuration
func WithMaxKeepaliveDuration(maxkt time.Duration) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.MaxKeepaliveDuration = maxkt
	}
}

// WithTCPKeepalive Whether the operating system should
// send tcp keep-alive messages on the tcp connection.
func WithTCPKeepalive() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.TCPKeepalive = true
	}
}

// WithTCPKeepalivePeriod set the time interval between tcp keep alive messages.
func WithTCPKeepalivePeriod(t time.Duration) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.TCPKeepalivePeriod = t
	}
}

// WithMaxRequestBodySize set MaxRequestBodySize
func WithMaxRequestBodySize(size int) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.MaxRequestBodySize = size
	}
}

// WithReduceMemoryUsage Aggressive memory usage reduction is disabled by default.
func WithReduceMemoryUsage() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.ReduceMemoryUsage = true
	}
}

// WithGetOnly Rejects all non-GET requests if set to true.
func WithGetOnly() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.GetOnly = true
	}
}

// WithDisablePreParseMultipartForm Will not pre parse Multipart Form data if set to true.
func WithDisablePreParseMultipartForm() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.DisablePreParseMultipartForm = true
	}
}

// WithLogAllErrors Whether to record all errors(debug)
func WithLogAllErrors() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.LogAllErrors = true
	}
}

// WithDisableHeaderNamesNormalizing Header names are passed as-is
// without normalization if this option is set.
func WithDisableHeaderNamesNormalizing() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.DisableHeaderNamesNormalizing = true
	}
}

// WithSleepWhenConcurrencyLimitsExceeded default [when is 0]: don't sleep
// and accept new connections immidiatelly
func WithSleepWhenConcurrencyLimitsExceeded(t time.Duration) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.SleepWhenConcurrencyLimitsExceeded = t
	}
}

// WithNoDefaultServerHeader NoDefaultServerHeader, when set to true,
// causes the default Server header to be excluded from the Response.
func WithNoDefaultServerHeader() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.NoDefaultServerHeader = true
	}
}

// WithNoDefaultDate NoDefaultDate, when set to true, causes the
// default Date header to be excluded from the Response.
func WithNoDefaultDate() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.NoDefaultDate = true
	}
}

// WithNoDefaultContentType DefaultContentType, when set to true, causes the
// default ContentType header to be excluded from the Response.
func WithNoDefaultContentType() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.NoDefaultContentType = true
	}
}

// WithConnState set ConnState hook
func WithConnState(f func(net.Conn, fasthttp.ConnState)) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.ConnState = f
	}
}

// WithLogger Logger which is used by RequestCtx.Logger().
// By default standard logger from log package is used.
func WithLogger(logger fasthttp.Logger) SvrOption {
	return func(svr *fasthttp.Server) {
		svr.Logger = logger
	}
}

// WithKeepHijackedConns KeepHijackedConns is an opt-in disable of
// connection close by fasthttp after connections' HijackHandler returns.
func WithKeepHijackedConns() SvrOption {
	return func(svr *fasthttp.Server) {
		svr.KeepHijackedConns = true
	}
}
