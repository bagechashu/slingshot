package server

import (
	"fmt"

	"slingshot/app"
	"slingshot/config"

	"github.com/labstack/echo/v4"
	"github.com/liushuochen/gotable"
	"github.com/spf13/cobra"
)

var E *echo.Echo

func init() {
	E = echo.New()
	app.InitRoutes(E)
}

func Run(cmd *cobra.Command, args []string) {
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
