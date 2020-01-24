package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

// grantTablesCmd represents the grantTables command
var grantTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: parser.ShortDesc("grant privileges on tables"),
	Long: `Usage: ( grant | revoke ) [privileges] on tables [except <table_name>[,<table_name> ...];

Examples:
  grant privileges on tables;
  grant on tables except this_schema.this_table;
  revoke on tables;
  revoke privileges on tables except that_schema.that_table,other_schema.other_table`,
	Run: runGrantOrRevokeCommand,
}
