package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// beginCmd represents the begin command
var beginCmd = &cobra.Command{
	Use:   "begin",
	Short: "Begins a transaction",
	Long:  `Usage: ddsl begin [transaction]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("begin called")
	},
}

func init() {
	rootCmd.AddCommand(beginCmd)
}
