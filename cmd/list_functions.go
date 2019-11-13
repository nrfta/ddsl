package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listFunctions = &cobra.Command{
	Use:   "functions",
	Short: parser.ShortDesc("list functions"),
	Long: `Usage: list functions [ (in | except in) <schema_name>[,<schema_name>...]];

Examples:
  list functions
  list functions in foo_schema
  list functions except in foo_schema,bar_schema
`,
	Run: runListCmd,
}
