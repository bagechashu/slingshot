package user

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	ug := e.Group("/user")
	ug.POST("/login", login)
	ug.POST("/register", register)
	ug.POST("/", addUser)
	ug.GET("/:uid", getUser)
	ug.GET("/", getUsers)
	ug.DELETE("/:uid", delUser)
	ug.GET("/:uid/roles", getRolesOfUser)

	rg := e.Group("/role")
	rg.POST("/", addRole)
	rg.POST("/:rid", getRole)
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
