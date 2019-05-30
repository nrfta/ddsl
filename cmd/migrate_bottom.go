package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// bottomCmd represents the bottom command
var bottomCmd = &cobra.Command{
	Use:   "bottom",
	Short: parser.ShortDesc("migrate bottom"),
	Long: `Example:
  migrate bottom;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bottom called")
	},
}

func init() {
	migrateCmd.AddCommand(bottomCmd)
}
