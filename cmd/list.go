package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// createCmd represents the create command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: parser.ShortDesc("list"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listSchemas)
	listCmd.AddCommand(listForeignKeys)
	listCmd.AddCommand(listSchemaItems)
	listCmd.AddCommand(listTables)
	listCmd.AddCommand(listViews)
	listCmd.AddCommand(listFunctions)
	listCmd.AddCommand(listProcedures)
	listCmd.AddCommand(listTypes)

	viper.BindEnv("output_format")

	rootCmd.PersistentFlags().StringP("format", "o", "text", "output format for list command (default DDSL_OUTPUT_FORMAT=text). May be text, csv, or json.")
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))

}

func runListCmd(cmd *cobra.Command, args []string) {
	command := fmt.Sprintf("list %s", cmd.Use)
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

