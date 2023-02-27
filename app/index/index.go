package index

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IndexGroup(g *echo.Group) {
	g.GET("/", hello)
	g.GET("/:name", helloName)
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func helloName(c echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, "Hello, "+name+"!")
}
