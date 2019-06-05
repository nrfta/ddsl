package cmd

import (
	"github.com/neighborly/ddsl/parser"
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
	Run: runCreateOrDropCommand,
}
