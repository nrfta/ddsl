package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// tableCmd represents the table command
var tableCmd = &cobra.Command{
	Use:   "table",
	Short: parser.ShortDesc("create table"),
	Long: `Usage: ( create | drop ) table <table_name>[,<table_name> ...];

Examples:
  create table this_schema.this_table;
  drop table that_schema.that_table,other_schema.other_table;`,
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}
