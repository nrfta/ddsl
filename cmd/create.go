package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	CREATE = "create"
	DROP = "drop"
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


func createOrDrop(cmd *cobra.Command) string {
	for _, a := range os.Args {
		switch a {
		case CREATE:
			return CREATE
		case DROP:
			return DROP
		}
	}
	return ""
}

func constructCreateOrDropCommand(cmd *cobra.Command, args []string) string {
	corD := createOrDrop(cmd)
	if len(corD) == 0 {
		panic("use only for create or drop")
	}
	command := fmt.Sprintf("%s %s", corD, cmd.Use)
	if len(args) > 0 {
		command += " "
	}
	command += strings.Join(args, " ")
	return command
}