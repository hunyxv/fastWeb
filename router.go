package fastweb

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees                  map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash  bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath      bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS          bool

	// An optional http.Handler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS          HandlerFunc
	
	// Cached value of global (*) allowed methods
	globalAllowed          string

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound               HandlerFunc

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed       HandlerFunc

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler           HandlerFunc
}

func newRouter() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
}

// // GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handle)
// func (r *Router) GET(path string, handle HandlerFunc) {
// 	r.addRoute(fasthttp.MethodGet, path, handle)
// }

// // HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handle)
// func (r *Router) HEAD(path string, handle HandlerFunc) {
// 	r.addRoute(fasthttp.MethodHead, path, handle)
// }

// // OPTIONS is a shortcut for router.Handle(fasthttp.MethodOptions, path, handle)
// func (r *Router) OPTIONS(path string, handle HandlerFunc) {
// 	r.addRoute(fasthttp.MethodOptions, path, handle)
// }

// // POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handle)
// func (r *Router) POST(path string, handle HandlerFunc) {
// 	r.addRoute(fasthttp.MethodPost, path, handle)
// }

// // PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handle)
// func (r *Router) PUT(path string, handle HandlerFunc) {
// 	r.addRoute(fasthttp.MethodPut, path, handle)
// }

// // PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handle)
// func (r *Router) PATCH(path string, handle HandlerFunc) {
// 	r.addRoute(fasthttp.MethodPatch, path, handle)
// }

// // DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handle)
// func (r *Router) DELETE(path string, handle HandlerFunc) {
// 	r.addRoute(fasthttp.MethodDelete, path, handle)
// }

func (r *Router) addRoute(method, path string, handle HandlerFunc) {
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}
	root, ok := r.trees[method]
	if !ok {
		root = new(node)
		r.trees[method] = root

		r.globalAllowed = r.allowed("*", "")
	}

	root.addRoute(path, handle)
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.trees {
				if method == fasthttp.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == fasthttp.MethodOptions {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, fasthttp.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) ServeFiles(path, root string) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := fasthttp.FSHandler(root, func(prefix string) (count int) {
		for i := 0; i < len(prefix); i++ {
			if prefix[i] == '/' {
				count++
			}
		}
		return
	}(path[:len(path)-10]))

	r.addRoute(fasthttp.MethodGet, path, func(ctx Context) {
		fileServer(ctx.GetFctx())
	})
}

func (r *Router) recv(ctx Context) {
	if rcv := recover(); rcv != nil {
		if r.PanicHandler != nil {
			ctx.SetUserValue("PanicError", rcv)
			r.PanicHandler(ctx)
		} else {
			ctx.Error(
				fasthttp.StatusMessage(fasthttp.StatusInternalServerError), 
				fasthttp.StatusInternalServerError,
			)
		}
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (HandlerFunc, Params, bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

func (r *Router) Handle(ctx *context) {
	defer r.recv(ctx)

	path := ctx.Path()
	if root := r.trees[ctx.Method()]; root != nil {
		if handle, ps, tsr := root.getValue(path); handle != nil {
			ctx.SetURLParam(ps)
			handle(ctx)
			return
		} else if ctx.Method() != fasthttp.MethodConnect && path != "/" {
			code := 301
			if ctx.Method() != fasthttp.MethodGet {
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					ctx.SetPath(path[:len(path)-1])
				} else {
					ctx.SetPath(path + "/")
				}
				ctx.fctx.Redirect(ctx.Path(), code)
				return
			}

			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					ctx.SetPath(b2s(fixedPath))
					ctx.Redirect(ctx.Path(), code)
					return
				}
			}
		}
	}

	if ctx.Method() == fasthttp.MethodOptions {
		if allow := r.allowed(path, ctx.Method()); allow != "" {
			ctx.SetHeader("Allow", allow)
			if r.MethodNotAllowed != nil {
				r.MethodNotAllowed(ctx)
			} else {
				ctx.Error(
					fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed), 
					fasthttp.StatusMethodNotAllowed,
				)
			}
			return
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound(ctx)
	} else{
		ctx.NotFound()
	}
}
