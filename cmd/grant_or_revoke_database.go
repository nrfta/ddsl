package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// grantDatabaseCmd represents the grantOrRevokeDatabase command
var grantDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: parser.ShortDesc("grant privileges on database"),
	Long: `Usage: ( grant | revoke ) [privileges] on database;

Examples:
  grant privileges on database;
  revoke on database;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grantOrRevokeDatabase called")
	},
}
