package user

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	g := e.Group("/user")

	g.GET("/users", getUsers)
	g.GET("/user/:id", getUser)
	g.POST("/user", addUser)
	g.DELETE("/user/:id", delUser)
}
