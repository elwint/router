package router

import urlpath "path"

func join(prefix, path string) string {
	return urlpath.Join(prefix, urlpath.Clean(path))
}

// Group is a router group with a shared prefix and set of middlewares
type Group struct {
	router     *Router
	prefix     string
	middleware []Middleware
}

// Group creates a new router group with a shared prefix and set of middlewares
func (g *Group) Group(prefix string, middleware ...Middleware) *Group {
	return &Group{prefix: join(g.prefix, prefix), router: g.router, middleware: append(g.middleware, middleware...)}
}

// GET adds a GET route
func (g *Group) GET(path string, handle Handle, middleware ...Middleware) {
	g.router.GET(join(g.prefix, path), handle, append(g.middleware, middleware...)...)
}

// POST adds a POST route
func (g *Group) POST(path string, handle interface{}, middleware ...Middleware) {
	g.router.POST(join(g.prefix, path), handle, append(g.middleware, middleware...)...)
}

// DELETE adds a DELETE route
func (g *Group) DELETE(path string, handle Handle, middleware ...Middleware) {
	g.router.DELETE(join(g.prefix, path), handle, append(g.middleware, middleware...)...)
}

// PUT adds a PUT route
func (g *Group) PUT(path string, handle interface{}, middleware ...Middleware) {
	g.router.PUT(join(g.prefix, path), handle, append(g.middleware, middleware...)...)
}

// PATCH adds a PATCH route
func (g *Group) PATCH(path string, handle interface{}, middleware ...Middleware) {
	g.router.PATCH(join(g.prefix, path), handle, append(g.middleware, middleware...)...)
}

// HEAD adds a HEAD route
func (g *Group) HEAD(path string, handle Handle, middleware ...Middleware) {
	g.router.HEAD(join(g.prefix, path), handle, append(g.middleware, middleware...)...)
}

// OPTIONS adds a OPTIONS route
func (g *Group) OPTIONS(path string, handle Handle, middleware ...Middleware) {
	g.router.OPTIONS(join(g.prefix, path), handle, append(g.middleware, middleware...)...)
}
