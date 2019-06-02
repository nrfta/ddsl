package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// foreignKeysCmd represents the foreignKeys command
var foreignKeysCmd = &cobra.Command{
	Use:   "foreign-keys",
	Short: parser.ShortDesc("create foreign-keys"),
	Long: `Examples:
  ddsl create foreign-keys;
  ddsl drop foreign-keys;`,
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}

