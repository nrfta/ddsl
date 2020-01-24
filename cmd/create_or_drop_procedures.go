package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// proceduresCmd represents the procedures command
var proceduresCmd = &cobra.Command{
	Use:   "procedures",
	Short: parser.ShortDesc("create procedures"),
	Long: `Usage: ( create | drop ) procedures [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];

Examples:
  create procedures;
  create procedures in this_schema;
  create procedures except in that_schema;
  drop procedures;
  drop procedures that_schema,other_schema;
  drop procedures except other_schema;`,
	Run: runCreateOrDropCommand,
}
