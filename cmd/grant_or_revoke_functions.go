package cmd

import (
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

var grantFunctionsCmd = &cobra.Command{
	Use:   "functions",
	Short: parser.ShortDesc("grant privileges on functions"),
	Long: `Usage: ( grant | revoke ) [privileges] on functions [except <function_name>[,<function_name> ...];

Examples:
  grant privileges on functions;
  grant on functions except this_schema.this_function;
  revoke on functions;
  revoke privileges on functions except that_schema.that_function,other_schema.other_function`,
	Run: runGrantOrRevokeCommand,
}
