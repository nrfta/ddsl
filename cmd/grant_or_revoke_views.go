package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// grantViewsCmd represents the grantViews command
var grantViewsCmd = &cobra.Command{
	Use:   "views",
	Short: parser.ShortDesc("grant privileges on views"),
	Long: `Usage: ( grant | revoke ) [privileges] on views [except <view_name>[,<view_name> ...];

Examples:
  grant privileges on views;
  grant on views except this_schema.this_view;
  revoke on views;
  revoke privileges on views except that_schema.that_view,other_schema.other_view`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grantViews called")
	},
}
