package router

import (
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

type route struct {
	Method     string
	Path       string
	Handle     interface{}
	Middleware []Middleware
}

// Handle handles a request
type Handle = func(*Context) error

// ErrorHandle handles a request
type ErrorHandle func(*Context, interface{})

// Middleware is a function that runs before your route, it gets the next handler as a parameter
type Middleware func(Handle) Handle

// Binder reads input to dst, returns true is successful
type Reader func(c *Context, dst interface{}) bool

// Router is the router itself
type Router struct {
	routes                  []route
	Reader                  Reader
	Renderer                Renderer
	middleware              []Middleware
	NotFoundHandler         Handle
	MethodNotAllowedHandler Handle
	ErrorHandler            ErrorHandle
}

// New returns a new Router
func New() *Router {
	return &Router{Reader: defaultReader, NotFoundHandler: defaultNotFoundHandler, MethodNotAllowedHandler: defaultMethodNotAllowedHandler, ErrorHandler: defaultErrorHandler}
}

// Use adds a global middleware
func (r *Router) Use(m ...Middleware) {
	r.middleware = append(r.middleware, m...)
}

// Group creates a new router group with a shared prefix and set of middlewares
func (r *Router) Group(prefix string, middleware ...Middleware) *Group {
	return &Group{prefix: prefix, router: r, middleware: middleware}
}

// GET adds a GET route
func (r *Router) GET(path string, handle Handle, middleware ...Middleware) {
	r.routes = append(r.routes, route{`GET`, path, handle, middleware})
}

// POST adds a POST route
func (r *Router) POST(path string, handle interface{}, middleware ...Middleware) {
	checkInterfaceHandle(handle)
	r.routes = append(r.routes, route{`POST`, path, handle, middleware})
}

// DELETE adds a DELETE route
func (r *Router) DELETE(path string, handle Handle, middleware ...Middleware) {
	r.routes = append(r.routes, route{`DELETE`, path, handle, middleware})
}

// PUT adds a PUT route
func (r *Router) PUT(path string, handle interface{}, middleware ...Middleware) {
	checkInterfaceHandle(handle)
	r.routes = append(r.routes, route{`PUT`, path, handle, middleware})
}

// PATCH adds a PATCH route
func (r *Router) PATCH(path string, handle interface{}, middleware ...Middleware) {
	checkInterfaceHandle(handle)
	r.routes = append(r.routes, route{`PATCH`, path, handle, middleware})
}

// HEAD adds a HEAD route
func (r *Router) HEAD(path string, handle Handle, middleware ...Middleware) {
	r.routes = append(r.routes, route{`HEAD`, path, handle, middleware})
}

// OPTIONS adds a OPTIONS route
func (r *Router) OPTIONS(path string, handle Handle, middleware ...Middleware) {
	r.routes = append(r.routes, route{`OPTIONS`, path, handle, middleware})
}

// Start starts the web server and binds to the given address
func (r *Router) Start(addr string) error {
	httpr := r.getHttpr()

	return http.ListenAndServe(addr, httpr)
}

func (r *Router) getHttpr() *httprouter.Router {
	httpr := httprouter.New()

	for _, v := range r.routes {
		handle, ok := v.Handle.(Handle)
		if !ok {
			handle = handlePOST(r, v.Handle)
		}

		httpr.Handle(v.Method, v.Path, handleReq(r, handle, append(r.middleware, v.Middleware...)))
	}

	httpr.NotFound = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handleReq(r, r.NotFoundHandler, r.middleware)(res, req, nil)
	})

	httpr.MethodNotAllowed = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handleReq(r, r.MethodNotAllowedHandler, r.middleware)(res, req, nil)
	})

	httpr.PanicHandler = func(res http.ResponseWriter, req *http.Request, err interface{}) {
		c := newContext(r, res, req, nil)
		r.ErrorHandler(c, err)
	}

	return httpr
}

func handleErr(errHandler ErrorHandle, err interface{}) Handle {
	return func(c *Context) error {
		errHandler(c, err)
		return nil
	}
}

func checkInterfaceHandle(f interface{}) {
	if _, ok := f.(Handle); ok {
		return
	}

	rt := reflect.TypeOf(f)

	if rt.Kind() != reflect.Func {
		panic(`non-func handle`)
	}

	if rt.NumIn() != 2 {
		panic(`handle should take 2 arguments`)
	}

	if rt.NumOut() != 1 || rt.Out(0).Name() != `error` {
		panic(`handle should return only error`)
	}

	if rt.In(0) != reflect.TypeOf(&Context{}) {
		panic(`handle should accept Context as first argument`)
	}

	return
}

func handlePOST(r *Router, f interface{}) Handle {
	funcRv, inputRt := reflect.ValueOf(f), reflect.TypeOf(f).In(1)

	return func(c *Context) error {
		data := reflect.New(inputRt)

		if !r.Reader(c, data.Interface()) {
			c.Request.Body.Close()
			return nil
		}
		c.Request.Body.Close()

		out := funcRv.Call([]reflect.Value{reflect.ValueOf(c), data.Elem()})

		if out[0].IsNil() {
			return nil
		}
		return out[0].Interface().(error)
	}
}

func handleReq(r *Router, handle Handle, m []Middleware) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, param httprouter.Params) {
		c := newContext(r, res, req, param)

		f := handle
		for i := len(m) - 1; i >= 0; i-- { // TODO: 1,2,3 of 3,2,1
			f = m[i](f)
		}

		err := f(c)

		if err != nil {
			r.ErrorHandler(c, err)
		}
	}
}
