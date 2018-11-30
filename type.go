package router

import (
	"io"
)

// Renderer renders a template
type Renderer interface {
	Render(w io.Writer, template string, data interface{}, c *Context) error
}
