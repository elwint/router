package router

import urlpath "path"

type Group struct {
	router *Router
	prefix string
}

func (g *Group) Group(prefix string) *Group {
	return &Group{prefix: urlpath.Join(g.prefix, prefix), router: g.router}
}

// GET adds a GET route
func (g *Group) GET(path string, handle GetHandle) {
	g.router.GET(urlpath.Join(g.prefix, path), handle)
}

// POST adds a POST route
func (g *Group) POST(path string, handle interface{}) {
	g.router.POST(urlpath.Join(g.prefix, path), handle)
}

// DELETE adds a DELETE route
func (g *Group) DELETE(path string, handle GetHandle) {
	g.router.DELETE(urlpath.Join(g.prefix, path), handle)
}

// PUT adds a PUT route
func (g *Group) PUT(path string, handle interface{}) {
	checkInterfaceHandle(handle)
	g.router.PUT(urlpath.Join(g.prefix, path), handle)
}

// PATCH adds a PATCH route
func (g *Group) PATCH(path string, handle interface{}) {
	checkInterfaceHandle(handle)
	g.router.PATCH(urlpath.Join(g.prefix, path), handle)
}

// HEAD adds a HEAD route
func (g *Group) HEAD(path string, handle GetHandle) {
	g.router.HEAD(urlpath.Join(g.prefix, path), handle)
}

// OPTIONS adds a OPTIONS route
func (g *Group) OPTIONS(path string, handle GetHandle) {
	g.router.OPTIONS(urlpath.Join(g.prefix, path), handle)
}
