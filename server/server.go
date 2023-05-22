package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slingshot/app"
	"slingshot/app/user"
	"slingshot/db"
	"slingshot/templates"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/liushuochen/gotable"
	"github.com/spf13/cobra"
)

// type slingshot struct {
// 	Mode        string
// 	RuntimeRoot string
// 	Server      *echo.Echo
// 	Log         *zap.SugaredLogger
// 	Validate    *validator.Validate
// }
// var s = &slingshot{}

var E *echo.Echo

func setup() {
	db.InitMysql()
	user.InitRbac()

	E = echo.New()

	// E.Use(middleware.Logger())
	E.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status} ${latency_human} ${error} \n",
		Output: os.Stdout,
	}))

	E.Use(middleware.Recover())
	// TODO: csrf middleware

	E.Renderer = templates.HTMLRender
}

func Setup(cmd *cobra.Command, args []string) {
	setup()
}

func Run(cmd *cobra.Command, args []string) {
	app.InitRoutes(E)
	// E.Logger.Fatal(E.Start(fmt.Sprintf(":%d", config.Cfg.Server.Port)))

	// Setup
	E.Logger.SetLevel(log.DEBUG)

	// Start server
	go func() {
		if err := E.Start(":1323"); err != nil && err != http.ErrServerClosed {
			E.Logger.Fatalf("shutting down the server: %s", err)
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
