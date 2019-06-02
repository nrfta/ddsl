package cmd

import (
	"github.com/neighborly/ddsl/log"
	"github.com/neighborly/ddsl/parser"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		code, err := runCLICommand(constructCreateOrDropCommand(cmd, args))
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(code)
	},
}
