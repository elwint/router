package router

import (
	"encoding/json"
	"net/http"
)

func defaultNotFoundHandler(c *Context) error {
	return c.StatusText(http.StatusNotFound)
}

func defaultMethodNotAllowedHandler(c *Context) error {
	return c.StatusText(http.StatusMethodNotAllowed)
}

func defaultErrorHandler(c *Context, err interface{}) {
	_ = c.StatusText(http.StatusInternalServerError)
}

func defaultReader(c *Context, dst interface{}) (bool, error) {
	err := json.NewDecoder(c.Request.Body).Decode(dst)
	if err != nil {
		return false, c.StatusText(http.StatusBadRequest)
	}

	return true, nil
}
