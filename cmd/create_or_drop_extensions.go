package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// extensionsCmd represents the extensions command
var extensionsCmd = &cobra.Command{
	Use:   "extensions",
	Short: parser.ShortDesc("create extensions"),
	Long: `Usage: ( create | drop ) extensions;

Examples:
  create extensions;
  drop extensions;`,
	Run: runCreateOrDropCommand,
}
