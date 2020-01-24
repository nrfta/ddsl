package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listProcedures = &cobra.Command{
	Use:   "procedures",
	Short: parser.ShortDesc("list procedures"),
	Long: `Usage: list procedures [ (in | except in) <schema_name>[,<schema_name>...]];

Examples:
  list procedures
  list procedures in foo_schema
  list procedures except in foo_schema,bar_schema
`,
	Run: runListCmd,
}
