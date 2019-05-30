package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: parser.ShortDesc("rollback"),
	Long: `Usage: rollback [transaction];`,
	Run: rollback,
}

func rollback(cmd *cobra.Command, args []string) {
	fmt.Println("rollback called")
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
