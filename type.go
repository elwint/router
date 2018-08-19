package router

import (
	"io"
)

type Renderer interface {
	Render(w io.Writer, template string, data interface{}, c *Context) error
}
