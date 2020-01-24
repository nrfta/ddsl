package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listViews = &cobra.Command{
	Use:   "views",
	Short: parser.ShortDesc("list views"),
	Long: `Usage: list views [ (in | except in) <schema_name>[,<schema_name>...]];

Examples:
  list views
  list views in foo_schema
  list views except in foo_schema,bar_schema
`,
	Run: runListCmd,
}
