package user

import (
	"github.com/labstack/echo/v4"
)

// TODO: JWT auth
func InitRoutes(e *echo.Echo) {
	ug := e.Group("/user")

	ug.POST("/", addUser)
	ug.GET("/:uid", getUser)
	ug.GET("/", getUsers)
	ug.DELETE("/:uid", delUser)
	ug.GET("/:uid/roles", getRolesOfUser)

	rg := e.Group("/role")
	rg.POST("/", addRole)
	rg.GET("/", getRoles)
	rg.DELETE("/:rid", delRole)
	rg.POST("/:rid/users", addUsersForRole)
	rg.GET("/:rid/users", getUsersOfRole)
	rg.GET("/:rid/policy", getRolePolicy)

	pg := e.Group("/policy")
	r_pg := pg.Group("/role")
	r_pg.GET("/", getPolicys)
	r_pg.POST("/", addPolicyForRole)
	r_pg.DELETE("/", delPolicyForRole)
}
