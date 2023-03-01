package app

import (
	"slingshot/app/index"
	"slingshot/app/user"
	"slingshot/templates"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	templates.InitRoutes(e)
	index.InitRoutes(e)
	user.InitRoutes(e)
}
