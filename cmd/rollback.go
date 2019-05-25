package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the active transaction",
	Long: `Usage: rollback [transaction]`,
	Run: rollback,
}

func rollback(cmd *cobra.Command, args []string) {
	fmt.Println("rollback called")
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
