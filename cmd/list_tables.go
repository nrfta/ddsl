package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listTables = &cobra.Command{
	Use:   "tables",
	Short: parser.ShortDesc("list tables"),
	Long: `Usage: list tables [ (in | except in) <schema_name>[,<schema_name>...]];

Examples:
  list tables
  list tables in foo_schema
  list tables except in foo_schema,bar_schema
`,
	Run: runListCmd,
}
