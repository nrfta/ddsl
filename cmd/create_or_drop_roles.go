package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: parser.ShortDesc("create roles"),
	Long: `Use caution: roles may be shared across databases.

Examples:
  create roles;
  drop roles;`,
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}
