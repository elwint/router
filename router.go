package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

type route struct {
	Method string
	Path   string
	Handle interface{}
}

// GetHandle handles a request that doesn't receive a body
type GetHandle func(*Context) error

// Router is the router itself
type Router struct {
	routes   []route
	Renderer Renderer
}

// New returns a new Router
func New() *Router {
	return &Router{}
}

func (r *Router) Group(prefix string) *Group {
	return &Group{prefix: prefix, router: r}
}

// GET adds a GET route
func (r *Router) GET(path string, handle GetHandle) {
	r.routes = append(r.routes, route{`GET`, path, handle})
}

// POST adds a POST route
func (r *Router) POST(path string, handle interface{}) {
	checkInterfaceHandle(handle)
	r.routes = append(r.routes, route{`POST`, path, handle})
}

// DELETE adds a DELETE route
func (r *Router) DELETE(path string, handle GetHandle) {
	r.routes = append(r.routes, route{`DELETE`, path, handle})
}

// PUT adds a PUT route
func (r *Router) PUT(path string, handle interface{}) {
	checkInterfaceHandle(handle)
	r.routes = append(r.routes, route{`PUT`, path, handle})
}

// PATCH adds a PATCH route
func (r *Router) PATCH(path string, handle interface{}) {
	checkInterfaceHandle(handle)
	r.routes = append(r.routes, route{`PATCH`, path, handle})
}

// HEAD adds a HEAD route
func (r *Router) HEAD(path string, handle GetHandle) {
	r.routes = append(r.routes, route{`HEAD`, path, handle})
}

// OPTIONS adds a OPTIONS route
func (r *Router) OPTIONS(path string, handle GetHandle) {
	r.routes = append(r.routes, route{`OPTIONS`, path, handle})
}

// Start starts the web server and binds to the given address
func (r *Router) Start(addr string) error {
	httpr := r.getHttpr()

	return http.ListenAndServe(addr, httpr)
}

func (r *Router) getHttpr() *httprouter.Router {
	httpr := httprouter.New()

	for _, v := range r.routes {
		if handle, ok := v.Handle.(GetHandle); ok {
			httpr.Handle(v.Method, v.Path, handleGET(r, handle))
			continue
		}

		httpr.Handle(v.Method, v.Path, handlePOST(r, v.Handle))
	}

	return httpr
}

func checkInterfaceHandle(f interface{}) {
	if _, ok := f.(GetHandle); ok {
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

func handlePOST(r *Router, f interface{}) httprouter.Handle {
	funcRv, inputRt := reflect.ValueOf(f), reflect.TypeOf(f).In(1)

	return func(res http.ResponseWriter, req *http.Request, param httprouter.Params) {
		c := newContext(r, res, req, param)

		data := reflect.New(inputRt)
		{
			err := json.NewDecoder(req.Body).Decode(data.Interface())
			req.Body.Close()
			if err != nil {
				c.NoContent(400) // TODO: send info about error (BindError)
				return
			}
		}

		out := funcRv.Call([]reflect.Value{reflect.ValueOf(c), data.Elem()})
		err := out[0].Interface()
		_ = err
	}
}

func handleGET(r *Router, f GetHandle) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, param httprouter.Params) {
		c := newContext(r, res, req, param)

		err := f(c)

		fmt.Println(err)
	}
}
