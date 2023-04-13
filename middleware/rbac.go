package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"slingshot/config"
	"sync"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/labstack/echo/v4"
)

type Casbin struct {
	once     sync.Once
	adapter  *xormadapter.Adapter
	Enforcer *casbin.Enforcer
}

// NewRbac
var Rbac = new(Casbin)

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

			if ok, _ := Rbac.Enforcer.Enforce(sub, obj, act); !ok {
				return c.HTML(http.StatusForbidden, "no permission")
			}
			return next(c)
		}
	}
}

func InitRbac() {
	Rbac.once.Do(func() {
		var err error
		Rbac.adapter, err = xormadapter.NewAdapterWithTableName("mysql", config.Cfg.Database.DSN(), "casbin_rule", "sys_", true)
		if err != nil {
			fmt.Printf("adapter err: %v", err)
		}
		Rbac.Enforcer, err = casbin.NewEnforcer("middleware/rbac.conf", Rbac.adapter)
		if err != nil {
			fmt.Printf("Enforcer err: %v", err)
		}
	})
}
