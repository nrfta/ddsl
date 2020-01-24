package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

var grantProceduresCmd = &cobra.Command{
	Use:   "procedures",
	Short: parser.ShortDesc("grant privileges on procedures"),
	Long: `Usage: ( grant | revoke ) [privileges] on procedures [except <procedure_name>[,<procedure_name> ...];

Examples:
  grant privileges on procedures;
  grant on procedures except this_schema.this_procedure;
  revoke on procedures;
  revoke privileges on procedures except that_schema.that_procedure,other_schema.other_procedure`,
	Run: runGrantOrRevokeCommand,
}
