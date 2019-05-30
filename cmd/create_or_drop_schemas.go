package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

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
		fmt.Println("schemas called")
	},
}

