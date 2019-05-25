package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// beginCmd represents the begin command
var beginCmd = &cobra.Command{
	Use:   "begin",
	Short: "Begins a transaction",
	Long: `Usage: begin [transaction]

This command is useful when composing a DDSL script.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("begin called")
	},
}

func init() {
	rootCmd.AddCommand(beginCmd)
}
