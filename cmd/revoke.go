package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// revokeCmd represents the revoke command
var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: parser.ShortDesc("revoke"),
	Long: `Revoke privileges on database objects`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)
	revokeCmd.AddCommand(grantPrivilegesCmd)
}
