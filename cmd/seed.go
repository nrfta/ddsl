package cmd

import (
	"fmt"
	"github.com/nrfta/ddsl/log"
	"github.com/nrfta/ddsl/parser"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: parser.ShortDesc("seed"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}


func runSeedCommand(cmd *cobra.Command, args []string) {
	command := "seed " + cmd.Use
	if len(args) > 0 {
		command += " "
	}
	command += strings.Join(args, " ")
	code, err := runCLICommand(command)
	if err != nil {
		log.Error(err.Error())
	}
	os.Exit(code)
}
