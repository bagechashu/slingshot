package user

import (
	"github.com/labstack/echo/v4"
)

// TODO: JWT auth
func InitRoutes(e *echo.Echo) {
	ug := e.Group("/user")

	ug.POST("/", addUser)
	ug.GET("/:id", getUser)
	ug.GET("/", getUsers)
	ug.DELETE("/:id", delUser)
	ug.GET("/:id/roles", getUserRoles)

	rg := e.Group("/role")
	rg.POST("/", addRole)
	rg.GET("/", getRoles)
	rg.DELETE("/:id", delRole)
	rg.POST("/:id/users", addUsersForRole)
	rg.GET("/:id/users", getRoleUsers)
	rg.GET("/:id/policy", getRolePolicy)

	pg := e.Group("/policy")
	r_pg := pg.Group("/role")
	r_pg.GET("/", getPolicys)
	r_pg.POST("/", addPolicyForRole)
	r_pg.DELETE("/:id", delPolicyForRole)
}
