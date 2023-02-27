package cmd

import (
	"slingshot/config"

	"github.com/spf13/cobra"
)

var cfgFile string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config related",
	Long:  `config related`,
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	Long:  `show config`,
	Run:   config.Show,
}

func cmdInit() {
	config.InitConfig(cfgFile)
}

func init() {
	cobra.OnInitialize(cmdInit)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(showCmd)
}
