package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var listForeignKeys = &cobra.Command{
	Use:   "foreign-keys",
	Short: parser.ShortDesc("list foreign-keys"),
	Long: `Usage: list foreign-keys;
`,
	Run: runListCmd,
}
