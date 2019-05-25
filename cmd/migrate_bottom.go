package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// bottomCmd represents the bottom command
var bottomCmd = &cobra.Command{
	Use:   "bottom",
	Short: "Migrates the database to the earliest version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bottom called")
	},
}

func init() {
	migrateCmd.AddCommand(bottomCmd)
}
