package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var seedScriptFiles []string

// cmdCmd represents the cmd command
var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Seeds the database through a shell command",
	Long: `Examples:
  seed cmd -f ./scripts/seed1.sh -f ./scripts/seed2.sh
  seed cmd ` + "`shell script as text`" + `

Note that shell script is enclosed in backticks and can be multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cmd called")
	},
}

func init() {
	seedCmd.AddCommand(cmdCmd)

	cmdCmd.PersistentFlags().StringSliceVarP(&seedScriptFiles, "file", "f", nil, "file containing seed script, may be provided more than once")
}
