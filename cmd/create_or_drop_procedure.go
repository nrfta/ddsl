package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// procedureCmd represents the procedure command
var procedureCmd = &cobra.Command{
	Use:   "procedure",
	Short: parser.ShortDesc("create procedure"),
	Long: `Usage: ( create | drop ) procedure <procedure_name>[,<procedure_name> ...]

Examples:
  create procedure this_procedure;
  drop procedure that_procedure,other_procedure;`,
	Run: runCreateOrDropCommand,
}
