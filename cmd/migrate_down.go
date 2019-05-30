package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

var downNumVersions uint

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: parser.ShortDesc("migrate down"),
	Long: `Examples:
  migrate down -n 2;`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("down called")
	},
}

func init() {
	migrateCmd.AddCommand(downCmd)

	downCmd.PersistentFlags().UintVarP(&downNumVersions, "num-versions", "n", 0, "number of versions to migrate down")
}
