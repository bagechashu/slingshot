package user

import (
	"fmt"
	"log"
	"net/http"
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

// TODO: use JWT get user identity
func CasbinRBACMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			// skip casbin middleware
			if c.Get("skip_casbin") != nil {
				return err
			}

			log.Printf("==========casbin middleware==========\n")

			u := User{}
			if err := c.Bind(&u); err != nil {
				log.Printf("user bind in casbin err: %v", err)
			}

			if exist, err := u.GetByUid(); err != nil {
				return c.HTML(http.StatusInternalServerError, fmt.Sprintf("user get err: %v", err))
			} else if !exist {
				return c.HTML(http.StatusUnauthorized, "user not exist")
			}

			log.Printf("user: %v", u)
			sub := u.Rid

			// FIXME: sub can't be empty
			if len(sub) == 0 {
				return c.HTML(http.StatusUnauthorized, "user don't have role")
			}

			obj := c.Request().URL.RequestURI()
			act := c.Request().Method

			log.Printf("sub: %v, obj: %v, act: %v", sub, obj, act)

			// FIXME: can obj use regex?
			if ok, _ := Rbac.Enforcer.Enforce(sub, obj, act); !ok {
				return c.HTML(http.StatusForbidden, "no permission")
			}
			return err
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
