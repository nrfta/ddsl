package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

// grantViewCmd represents the grantView command
var grantViewCmd = &cobra.Command{
	Use:   "view",
	Short: parser.ShortDesc("grant privileges on view"),
	Long: `Usage: ( grant | revoke ) [privileges] on view <view_name>[,<view_name> ...];

Examples:
  grant privileges on view this_schema.this_view;
  revoke on view that_schema.that_view,other_schema.other_view;`,
	Run: runGrantOrRevokeCommand,
}
