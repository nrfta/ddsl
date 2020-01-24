package cmd

import (
"github.com/nrfta/ddsl/parser"
"github.com/spf13/cobra"
)

// triggersCmd represents the triggers command
var triggersCmd = &cobra.Command{
	Use:   "triggers",
	Short: parser.ShortDesc("create triggers"),
	Long: `Usage: ( create | drop ) triggers [on] <table_name>[,<table_name> ...];

Examples:
  create triggers on this_schema.this_table;
  drop triggers that_schema.that_table,other_schema.other_table;`,
	Run: runCreateOrDropCommand,
}
