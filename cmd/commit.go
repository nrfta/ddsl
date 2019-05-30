package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: parser.ShortDesc("commit"),
	Long: `Usage: commit [transaction];`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("commit called")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
