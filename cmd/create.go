package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
	"os"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: parser.ShortDesc("create"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	addCreateOrDropSubCmds(createCmd)
}
