package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// grantSchemaCmd represents the grantSchema command
var grantSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: parser.ShortDesc("grant privileges on schema"),
	Long: `Usage: ( grant | revoke ) [privileges] on schema <schema_name>[,<schema_name> ...];

Examples:
  grant privileges on schema this_schema;
  revoke on schema that_scheme,other_schema;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grantSchema called")
	},
}

