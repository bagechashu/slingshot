package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slingshot/app"
	"slingshot/db"
	mw "slingshot/middleware"
	"slingshot/templates"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/liushuochen/gotable"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type slingshot struct {
	Mode        string
	RuntimeRoot string
	Server      *echo.Echo
	Log         *zap.SugaredLogger
	Validate    *validator.Validate
}

var E *echo.Echo

var s = &slingshot{}

func setup() {
	db.InitMysql()
	mw.InitRbac()

	mw.Rbac.Enforcer.AddPolicy("slingshot", "user", "write")
	s.Mode = "debug"
	// e.LoadPolicy()
	E = echo.New()
	E.Use(middleware.Logger())
	E.Use(middleware.Recover())

	// E.Use(mw.CheckPermission())
	E.Renderer = templates.HTMLRender
}

func Setup(cmd *cobra.Command, args []string) {
	setup()
}

func Run(cmd *cobra.Command, args []string) {
	app.InitRoutes(E)
	// E.Logger.Fatal(E.Start(fmt.Sprintf(":%d", config.Cfg.Server.Port)))

	// Setup
	E.Logger.SetLevel(log.INFO)

	// Start server
	go func() {
		if err := E.Start(":1323"); err != nil && err != http.ErrServerClosed {
			E.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := E.Shutdown(ctx); err != nil {
		E.Logger.Fatal(err)
	}
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
