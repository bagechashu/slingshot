package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os/user"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/labstack/echo/v4"
)

var (
	casbinAdapter  *xormadapter.Adapter
	CasbinEnforcer *casbin.Enforcer
)

func CheckPermission() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			obj := c.Request().URL.RequestURI()
			act := c.Request().Method
			var u user.User
			log.Printf("========================")
			if err := c.Bind(&u); err != nil {
				log.Printf("user bind in casbin err: %v", err)
			}
			sub := u.Username
			log.Printf("sub: %v, obj: %v, act: %v", sub, obj, act)
			log.Printf("========================")

			if ok, _ := CasbinEnforcer.Enforce(sub, obj, act); !ok {
				return c.HTML(http.StatusForbidden, "no permission")
			}
			return next(c)
		}
	}
}

func init() {
	var err error
	casbinAdapter, err = xormadapter.NewAdapter("mysql", "root:123456@tcp(127.0.0.1:3306)/demo?charset=utf8mb4", true)
	if err != nil {
		fmt.Printf("casbinAdapter err: %v", err)
	}
	CasbinEnforcer, err = casbin.NewEnforcer("./rbac_models.conf", casbinAdapter)
	if err != nil {
		fmt.Printf("CasbinEnforcer err: %v", err)
	}
}
