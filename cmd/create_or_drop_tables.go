package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// tablesCmd represents the tables command
var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: parser.ShortDesc("create tables"),
	Long: `Usage: ( create | drop ) tables [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];

Examples:
  create tables;
  create tables in this_schema;
  create tables except in that_schema;
  drop tables;
  drop tables that_schema,other_schema;
  drop tables except other_schema;`,
	Run: runCreateOrDropCommand,
}
