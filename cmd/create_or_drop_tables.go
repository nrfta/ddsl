package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}
