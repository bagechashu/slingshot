package app

import (
	"slingshot/app/user"

	"github.com/spf13/cobra"
)

func AppMigrate(cmd *cobra.Command, args []string) {
	// Migrate your models here
	user.Migrate()
}
