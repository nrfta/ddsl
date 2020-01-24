package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

var seedScriptFiles []string

// cmdCmd represents the cmd command
var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: parser.ShortDesc("seed cmd"),
	Long: `Examples:
  seed cmd -f ./scripts/seed1.sh -f ./scripts/seed2.sh;
  seed cmd ` + "`shell script as text`;" + `

Note that shell script is enclosed in backticks and can be multiple lines`,
	Run: runSeedCommand,
}

func init() {
	seedCmd.AddCommand(cmdCmd)

	cmdCmd.PersistentFlags().StringSliceVarP(&seedScriptFiles, "file", "f", nil, "file containing seed script, may be provided more than once")
}
