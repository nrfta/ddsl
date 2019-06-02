package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: parser.ShortDesc("migrate"),
	Long: `Examples:
  migrate up 1;
  migrate down 2;
  migrate top;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
