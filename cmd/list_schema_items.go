package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listSchemaItems = &cobra.Command{
	Use:   "schema-items",
	Short: parser.ShortDesc("list schema-items"),
	Long: `Usage: list schema-items [ (in | except in) <schema_name>[,<schema_name>...]];

Examples:
  list schema-items
  list schema-items in foo_schema
  list schema-items except in foo_schema,bar_schema
`,
	Run: runListCmd,
}
