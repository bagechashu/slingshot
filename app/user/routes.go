package user

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {

	e.POST("/login", login)
	e.POST("/register", register)

	ug := e.Group("/user", UserAuthMiddleware(SkipPathNoLimit))
	{
		ug.POST("/", addUser)
		ug.GET("/:uid", getUser)
		ug.GET("/", getUsers)
		ug.DELETE("/:uid", delUser)
		ug.GET("/:uid/roles", getRolesOfUser)
	}

	rg := e.Group("/role", UserAuthMiddleware(SkipPathNoLimit))
	{
		rg.POST("/", addRole)
		rg.POST("/:rid", getRole)
		rg.GET("/", getRoles)
		rg.DELETE("/:rid", delRole)
		rg.POST("/:rid/users", addUsersForRole)
		rg.GET("/:rid/users", getUsersOfRole)
		rg.GET("/:rid/policy", getRolePolicy)
	}

	pg := e.Group("/policy", UserAuthMiddleware(SkipPathNoLimit))
	{
		r_pg := pg.Group("/role")
		r_pg.GET("/", getPolicys)
		r_pg.POST("/", addPolicyForRole)
		r_pg.DELETE("/", delPolicyForRole)
	}
}
