package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listSchemas = &cobra.Command{
	Use:   "schemas",
	Short: parser.ShortDesc("list schemas"),
	Long: `Usage: list schemas;
`,
	Run: runListCmd,
}
