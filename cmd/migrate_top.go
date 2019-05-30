package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: parser.ShortDesc("migrate top"),
	Long: `Example:
  migrate top;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("top called")
	},
}

func init() {
	migrateCmd.AddCommand(topCmd)
}
