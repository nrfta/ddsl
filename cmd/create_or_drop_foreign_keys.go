package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// foreignKeysCmd represents the foreignKeys command
var foreignKeysCmd = &cobra.Command{
	Use:   "foreign-keys",
	Short: parser.ShortDesc("create foreign-keys"),
	Long: `Usage: ( create | drop ) foreign-keys [on <table_name>[,<table_name> ...]];
       (create | drop ) foreign-keys [in <schema_name>[,<schema_name> ...]];

Examples:
  create foreign-keys;
  create foreign-keys in this_schema;
  drop foreign-keys except in this_schema,that_schema;
  create foreign-keys on this_schema.this_table;
  drop foreign-keys on that_schema.that_table,other_schema.other_table;`,
	Run: runCreateOrDropCommand,
}

