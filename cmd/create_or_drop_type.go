package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// typeCmd represents the types command
var typeCmd = &cobra.Command{
	Use:   "type",
	Short: parser.ShortDesc("create type"),
	Long: `Usage: ( create | drop ) type  <type_name>[,<type-name> ...];

Examples:
  create type this_schema.this_type;
  drop types that_schema.that_type,other_schema.other_type;`,
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}
