package user

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	g := e.Group("/user")

	g.GET("/:id", getUser)
	g.POST("/", addUser)
	g.DELETE("/:id", delUser)
	g.POST("/grant/:id", addPolicy)
}
