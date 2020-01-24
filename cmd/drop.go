package cmd

import (
	"fmt"
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
	"os"
)

// dropCmd represents the drop command
var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: parser.ShortDesc("drop"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(dropCmd)
	addCreateOrDropSubCmds(dropCmd)
}
