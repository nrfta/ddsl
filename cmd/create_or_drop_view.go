package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: parser.ShortDesc("create view"),
	Long: `Usage: ( create | drop ) view <view_name>[,<view_name> ...];

Examples:
  create view this_schema.this_view;
  drop view that_schema.that_view,other_schema.other_view;`,
	Run: runCreateOrDropCommand,
}
