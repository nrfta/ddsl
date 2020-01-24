package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// constraintsCmd represents the constraints command
var constraintsCmd = &cobra.Command{
	Use:   "constraints",
	Short: parser.ShortDesc("create constraints"),
	Long: `Usage: ( create | drop ) constraints [on] <table_name>[,<table_name> ...];

Examples:
  create constraints on this_schema.this_table;
  drop constraints that_schema.that_table;`,
	Run: runCreateOrDropCommand,
}
