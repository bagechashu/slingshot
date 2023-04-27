package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"slingshot/config"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/labstack/echo/v4"
)

type Casbin struct {
	once     sync.Once
	adapter  *xormadapter.Adapter
	Enforcer *casbin.SyncedEnforcer
}

// NewRbac
var Rbac = new(Casbin)

func CasbinRBACMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var u user.User
			if err := c.Bind(&u); err != nil {
				log.Printf("user bind in casbin err: %v", err)
			}
			sub := u.Username

			obj := c.Request().URL.RequestURI()
			act := c.Request().Method

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
		Rbac.adapter, err = xormadapter.NewAdapterWithTableName(
			"mysql",
			config.Cfg.Database.DSN(),
			"casbin_rule",
			"sys_",
			true)
		if err != nil {
			fmt.Printf("adapter err: %v", err)
		}

		casbinModel, err := model.NewModelFromString(config.RBAC_MODEL)
		if err != nil {
			fmt.Printf("model err: %v", err)
		}
		Rbac.Enforcer, err = casbin.NewSyncedEnforcer(casbinModel, Rbac.adapter)

		// Rbac.Enforcer, err = casbin.NewSyncedEnforcer("middleware/rbac.conf", Rbac.adapter)
		if err != nil {
			fmt.Printf("Enforcer err: %v", err)
		}
	})
}
