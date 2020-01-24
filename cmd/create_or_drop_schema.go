package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: parser.ShortDesc("create schema"),
	Long: `Usage: ( create | drop ) schema <schema_name>[,<schema_name> ...]

Examples:
  create schema this_schema;
  drop schema that_schema,other_schema;`,
	Run: runCreateOrDropCommand,
}
