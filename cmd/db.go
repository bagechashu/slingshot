package cmd

import (
	"slingshot/app"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "db",
	Long:  `db`,
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Long:  `migrate`,
	Run:   app.AppMigrate,
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(migrateCmd)
}
