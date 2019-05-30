package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// typesCmd represents the types command
var typesCmd = &cobra.Command{
	Use:   "types",
	Short: parser.ShortDesc("create types"),
	Long: `Usage: ( create | drop ) types [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];

Examples:
  create types;
  create types in this_schema;
  create types except in that_schema;
  drop types;
  drop types that_schema,other_schema;
  drop types except other_schema;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("types called")
	},
}
