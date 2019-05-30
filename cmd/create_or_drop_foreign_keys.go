package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// foreignKeysCmd represents the foreignKeys command
var foreignKeysCmd = &cobra.Command{
	Use:   "foreign-keys",
	Short: parser.ShortDesc("create foreign-keys"),
	Long: `Examples:
  ddsl create foreign-keys;
  ddsl drop foreign-keys;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("foreign-keys called")
	},
}

