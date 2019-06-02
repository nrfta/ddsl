package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}
