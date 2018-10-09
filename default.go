package router

import (
	"fmt"
)

func defaultNotFoundHandler(c *Context) error {
	return c.String(404, `not found`)
}

func defaultMethodNotAllowedHandler(c *Context) error {
	return c.String(504, `method not allowed`)
}

func defaultErrorHandler(c *Context, err interface{}) {
	fmt.Println(err)
	c.String(500, `internal server error`)
}
