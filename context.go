package router

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Context is passed to handlers and middlewares
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Param    func(string) string
	store    map[string]interface{}
}

func newContext(res http.ResponseWriter, req *http.Request, param httprouter.Params) *Context {
	return &Context{req, res, param.ByName, make(map[string]interface{})}
}

// String returns the given status code and writes the bytes to the body
func (c *Context) Bytes(code int, b []byte) error {
	c.Response.WriteHeader(code)
	_, err := c.Response.Write(b)
	return err
}

// String returns the given status code and writes the string to the body
func (c *Context) String(code int, s string) error {
	c.Response.WriteHeader(code)
	_, err := c.Response.Write([]byte(s))
	return err
}

// NoContent returns the given status code without writing anything to the body
func (c *Context) NoContent(code int) error {
	c.Response.WriteHeader(code)
	return nil
}

// JSON returns the given status code and writes JSON to the body
func (c *Context) JSON(code int, data interface{}) error {
	c.Response.WriteHeader(code)
	return json.NewEncoder(c.Response).Encode(data) // TODO: Encode to buffer first to prevent partial responses on error
}

// Set sets a value in the context. Set is not safe to be used concurrently
func (c *Context) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get retrieves a value from the context.
func (c *Context) Get(key string) interface{} {
	return c.store[key]
}
