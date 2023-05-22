package templates

import (
	"embed"
	"slingshot/config"

	"github.com/labstack/echo/v4"
)

//go:embed assets
var assetsFS embed.FS

func InitRoutes(e *echo.Echo) {
	if config.Cfg.Server.Debug {
		e.Static("/assets", "templates/assets")
		return
	}
	e.StaticFS("/", assetsFS)
}
