package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

// grantTableCmd represents the grantTable command
var grantTableCmd = &cobra.Command{
	Use:   "table",
	Short: parser.ShortDesc("grant privileges on table"),
	Long: `Usage: ( grant | revoke ) [privileges] on table <table_name>[,<table_name> ...];

Examples:
  grant privileges on table this_schema.this_table;
  revoke on table that_schema.that_table,other_schema.other_table;`,
	Run: runGrantOrRevokeCommand,
}
