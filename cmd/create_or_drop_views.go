package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// viewsCmd represents the views command
var viewsCmd = &cobra.Command{
	Use:   "views",
	Short: parser.ShortDesc("create views"),
	Long: `Usage: ( create | drop ) views [[ ( in | except [in] ) ] <schema_name>[,<exclude_view_name> ...]];

Examples:
  create views;
  create views in this_schema;
  create views this_schema;
  drop views in that_schema,other_schema;
  drop views except in that_schema;
  drop views except that_schema,other_schema;`,
	Run: runCreateOrDropCommand,
}

