package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slingshot/app"
	"slingshot/db"
	"slingshot/middleware"
	"slingshot/templates"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"
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

var s = &slingshot{}

func globalVarInit() {
	db.InitMysql()
	middleware.RBAC.Init()

	s.Server = echo.New()
	s.Server.Renderer = templates.HTMLRender
}

func registerMiddleware() {

	s.Server.Use(middleware.RBAC.Middleware())

	app.InitRoutes(s.Server)
}

func Run(cmd *cobra.Command, args []string) {
	globalVarInit()
	registerMiddleware()
	// E.Logger.Fatal(E.Start(fmt.Sprintf(":%d", config.Cfg.Server.Port)))

	// Setup
	s.Server.Logger.SetLevel(log.INFO)

	// Start server
	go func() {
		if err := s.Server.Start(":1323"); err != nil && err != http.ErrServerClosed {
			s.Server.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		s.Server.Logger.Fatal(err)
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

	for _, r := range s.Server.Routes() {
		table.AddRow([]string{r.Method, r.Path, r.Name})
		// fmt.Printf("%+v", row)
	}

	fmt.Println(table)

}
