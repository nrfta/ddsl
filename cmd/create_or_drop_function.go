package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// functionCmd represents the function command
var functionCmd = &cobra.Command{
	Use:   "function",
	Short: parser.ShortDesc("create function"),
	Long: `Usage: ( create | drop ) function <function_name>[,<function_name> ...]

Examples:
  create function this_function;
  drop function that_function,other_function;`,
	Run: runCreateOrDropCommand,
}
