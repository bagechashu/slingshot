package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slingshot",
	Short: "A simple Go Web Frame for Beginner and Amateur",
	Long: `A simple Go Web Frame for Beginner and Amateur,
Hope it can help you.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
