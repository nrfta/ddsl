package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

var grantFunctionCmd = &cobra.Command{
	Use:   "function",
	Short: parser.ShortDesc("grant privileges on function"),
	Long: `Usage: ( grant | revoke ) [privileges] on function <function_name>[,<function_name> ...];

Examples:
  grant privileges on function this_schema.this_function;
  revoke on function that_schema.that_function,other_schema.other_function;`,
	Run: runGrantOrRevokeCommand,
}
