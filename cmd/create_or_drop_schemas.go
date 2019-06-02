package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var schemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: parser.ShortDesc("create schemas"),
	Long: `Usage: ( create | drop ) schemas [[except] <exclude_schema_name>[,<exclude_schema_name ...]];

Examples:
  create schemas;
  drop schemas except not_this_schema,nor_that_schema;`,
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}

