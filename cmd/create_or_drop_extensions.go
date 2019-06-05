package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// extensionsCmd represents the extensions command
var extensionsCmd = &cobra.Command{
	Use:   "extensions",
	Short: parser.ShortDesc("create extensions"),
	Long: `Usage: ( create | drop ) extensions [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];

Examples:
  create extensions;
  create extensions in this_schema;
  create extensions except in that_schema;
  drop extensions;
  drop extensions that_schema,other_schema;
  drop extensions except other_schema;`,
	Run: runCreateOrDropCommand,
}
