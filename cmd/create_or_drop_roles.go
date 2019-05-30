package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: parser.ShortDesc("create roles"),
	Long: `Use caution: roles may be shared across databases.

Examples:
  create roles;
  drop roles;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("roles called")
	},
}
