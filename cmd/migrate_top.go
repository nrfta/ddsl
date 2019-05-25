package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Migrates the database to the latest version",
	Long: `Example:
  migrate top`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("top called")
	},
}

func init() {
	migrateCmd.AddCommand(topCmd)
}
