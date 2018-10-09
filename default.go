package router

import (
	"encoding/json"
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

func defaultReader(c *Context, dst interface{}) bool {
	err := json.NewDecoder(c.Request.Body).Decode(dst)
	if err != nil {
		c.NoContent(400)
		return false
	}

	return true
}
