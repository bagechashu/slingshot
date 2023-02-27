package app

import (
	"slingshot/app/index"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	index.IndexGroup(e.Group(""))
}
