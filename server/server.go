package server

import (
	"fmt"
	"log"

	"slingshot/app"
	"slingshot/config"
	"slingshot/middleware"
	"slingshot/templates"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	casbin_mw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
	"github.com/liushuochen/gotable"
	"github.com/spf13/cobra"
)

var (
	E  *echo.Echo
	GA *gormadapter.Adapter
	CE *casbin.Enforcer
)

func registerMiddleware() {
	GA, err := middleware.NewAdapter()
	if err != nil {
		log.Fatalf("failed to create new adapter: %v", err)
	} else if GA == nil {
		log.Fatalf("failed to create new adapter: %v", err)
	}
	// log.Printf("adapter: %+v", GA)
	// GA.AddPolicy("9", "/api/v1/roles", []string{"GET", "POST", "PUT", "DELETE"})

	CE, err := middleware.NewEnforcer(config.RBAC_MODEL, *GA)
	if err != nil {
		log.Fatalf("failed to create new enforcer: %v", err)
	}

	E = echo.New()
	E.Renderer = templates.HTMLRender
	E.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		Enforcer: CE,
	}))
	app.InitRoutes(E)
}

func Run(cmd *cobra.Command, args []string) {
	registerMiddleware()
	E.Logger.Fatal(E.Start(fmt.Sprintf(":%d", config.Cfg.Server.Port)))
}

func WalkRoutes(cmd *cobra.Command, args []string) {
	table, err := gotable.Create(
		"Method",
		"Path",
		"Name",
	)
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return
	}

	table.Align("Method", 1)
	table.Align("Path", 1)
	table.Align("Name", 1)

	for _, r := range E.Routes() {
		table.AddRow([]string{r.Method, r.Path, r.Name})
		// fmt.Printf("%+v", row)
	}

	fmt.Println(table)

}
