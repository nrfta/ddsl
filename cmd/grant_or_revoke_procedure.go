package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

var grantProcedureCmd = &cobra.Command{
	Use:   "procedure",
	Short: parser.ShortDesc("grant privileges on procedure"),
	Long: `Usage: ( grant | revoke ) [privileges] on procedure <procedure_name>[,<procedure_name> ...];

Examples:
  grant privileges on procedure this_schema.this_procedure;
  revoke on procedure that_schema.that_procedure,other_schema.other_procedure;`,
	Run: runGrantOrRevokeCommand,
}
