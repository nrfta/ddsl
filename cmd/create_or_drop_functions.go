package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// functionsCmd represents the functions command
var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: parser.ShortDesc("create functions"),
	Long: `Usage: ( create | drop ) functions [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];

Examples:
  create functions;
  create functions in this_schema;
  create functions except in that_schema;
  drop functions;
  drop functions that_schema,other_schema;
  drop functions except other_schema;`,
	Run: runCreateOrDropCommand,
}
