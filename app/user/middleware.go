package user

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func UserAuthMiddleware(skipper middleware.Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			log.Printf("========== user auth middleware ==========\n")

			CheckJwtToken(c)
			RbacVarifiy(c)

			return next(c)
		}
	}
}

func SkipPathNoLimit(c echo.Context) bool {
	path := c.Path()
	pathsNoLimit, err := GetPathNoLimit()
	if err != nil {
		return false
	}

	// static file req path is "/*"
	for _, v := range pathsNoLimit {
		if path == v {
			log.Printf("=========== skip path: %v ============\n", path)
			return true
		}
	}

	return false
}
