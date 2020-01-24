package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

// grantSchemasCmd represents the grantSchemas command
var grantSchemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: parser.ShortDesc("grant privileges on schemas"),
	Long: `Usage: ( grant | revoke ) [privileges] on schemas [except <schema_name>[,<schema_name> ...]];

Examples:
  grant privileges on schemas;
  grant on schemas except this_schema;
  revoke privileges on schemas except that_schema,other_schema;`,
	Run: runGrantOrRevokeCommand,
}
