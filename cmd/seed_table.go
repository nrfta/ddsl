package cmd

import (
	"github.com/nrfta/ddsl/parser"
	"github.com/spf13/cobra"
)

// seedTableCmd represents the seedTable command
var seedTableCmd = &cobra.Command{
	Use:   "table",
	Short: parser.ShortDesc("seed table"),
	Long: `Usage: seed table <table_name>[,<table_name> ...];

This command executes seed scripts in the following location:
- ./<schema_name>/tables/<table_name>.seed.*

If the file extension is "sql" then the script is run directly on the database, otherwise
it is run as a shell script with "sh <seed_file>".

Examples:
  seed table this_schema.this_table;
  seed table this_schema.this_table,that_schema.that_table;`,
	Run: runSeedCommand,
}

func init() {
	seedCmd.AddCommand(seedTableCmd)
}
