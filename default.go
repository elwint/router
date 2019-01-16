package router

import (
	"encoding/json"
	"log"
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

func defaultPanicHandler(c *Context, v interface{}) {
	log.Println(c.Request.Method, c.Request.URL.Path+`: panic: `, v)

	defaultErrorHandler(c, v)
}

func defaultReader(c *Context, dst interface{}) (bool, error) {
	err := json.NewDecoder(c.Request.Body).Decode(dst)
	if err != nil {
		return false, c.StatusText(http.StatusBadRequest)
	}

	return true, nil
}
