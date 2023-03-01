package index

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	g := e.Group("")

	g.GET("/", index)
}
