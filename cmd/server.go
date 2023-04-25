package cmd

import (
	"slingshot/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  `server`,
}

var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "start server",
	Long:   `start server`,
	PreRun: server.Setup,
	Run:    server.Run,
}

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "walk route server",
	Long:  `walk route server`,
	Run:   server.WalkRoutes,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(runCmd)
	serverCmd.AddCommand(walkCmd)
}
