package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listTypes = &cobra.Command{
	Use:   "types",
	Short: parser.ShortDesc("list types"),
	Long: `Usage: list types [ (in | except in) <schema_name>[,<schema_name>...]];

Examples:
  list types
  list types in foo_schema
  list types except in foo_schema,bar_schema
`,
	Run: runListCmd,
}
