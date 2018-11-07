# router

Router is a HTTP router for Golang. It's built on [httprouter](https://github.com/julienschmidt/httprouter) and takes inspiration from [labstack/echo](https://github.com/labstack/echo), but with reduced complexity and easier data binding.

The data binding is made easier by specifying your input as a parameter to your function

Example:

```golang
type someType struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func handlePOST(c *router.Context, input someType) error {
	fmt.Println(input)
	return c.NoContent(200)
}
```

### Why make data binding shorter?

Many applications read, bind and validate data for most calls. In Echo this could mean adding boilerplate code to every call. This extra boilerplate code can make your code significantly longer and very hard to read.

```golang
func handlePOST(c *echo.Context) error {
	var input someType
	err := c.Bind(&input)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err = c.Validate(input)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Your actual code
}
```

In router you define your code to read, bind and validate the input once, and it applies to every POST, PATCH and PUT call.

### What about the performance of dynamic parameters?

While using dynamic parameters takes a bit of extra processing, this barely has any impact on performance. Router can still handle tens to hunderds of thousands of requests per second.

### How does middleware work?

Middleware works similar to most other routers. The dynamic parameters has no effect on middleware, input data is parsed after all middlewares and right before your handler.

## Installation

```bash
go get git.fuyu.moe/Fuyu/router
```

## Getting started

```golang
package main

import "git.fuyu.moe/Fuyu/router"

func main() {
	// Create a router instance
	r := router.New()
	
	// Add routes
	r.GET(`/`, yourGetFunc)
	r.POST(`/`, yourPostFunc)
	
	// Start router
	panic(r.Start(`127.0.0.1:8080`))
}
```

## Advice

### Configuration

For a serious project you should set `r.Reader`, `r.ErrorHandler`, `r.NotFoundHandler`, and `r.MethodNotAllowedHandler`.

### Templating

You can set `r.Renderer` and call `c.Render(code, tmplName, data)` in your handlers



## Examples

```golang
package main

import (
	"fmt"
	"time"

	"git.fuyu.moe/Fuyu/router"
)

type product struct {
	Name string `json:"name"`
}

func main() {
	r := router.New()
	r.Use(accessLog)

	r.GET(`/hello`, hello)

	a := r.Group(`/api`, gandalf)
	a.POST(`/product`, createProduct)
	a.PATCH(`/product/:id`, updateProduct)
	a.POST(`/logout`, logout)

	r.Start(`:8080`)
}

func accessLog(next router.Handle) router.Handle {
	return func(c *router.Context) error {
		t := time.Now()
		err := next(c)

		fmt.Println(c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, t, time.Since(t))

		return err
	}
}

func gandalf(next router.Handle) router.Handle {
	youShallPass := false
	return func(c *router.Context) error {
		if !youShallPass {
			return c.String(401, `You shall not pass`)
		}

		return next(c)
	}
}

func hello(c *router.Context) error {
	return c.String(200, `Hello`)
}

func createProduct(c *router.Context, p product) error {
	return c.JSON(200, p)
}

func updateProduct(c *router.Context, p product) error {
	productID := c.Param(`id`)
	return c.String(200, fmt.Sprintf(
		`ProductID %d new name %s`, productID, p.Name,
	))
}

func logout(c *router.Context) error {
	return c.String(200, `logout`)
}
```
