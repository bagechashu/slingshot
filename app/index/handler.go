package index

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	d := map[string]interface{}{
		"title": "Hello, World!",
	}
	return c.Render(http.StatusOK, "index/index.html", d)
}
