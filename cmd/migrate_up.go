package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

var upNumVersions uint

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: parser.ShortDesc("migrate up"),
	Long: `Examples:
  migrate up -n 2;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("up called")
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)

	upCmd.PersistentFlags().UintVarP(&upNumVersions, "num-versions", "n", 0, "number of versions to migrate up")
}
