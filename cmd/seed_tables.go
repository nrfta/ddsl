package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// seedTablesCmd represents the seedTables command
var seedTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: parser.ShortDesc("seed tables"),
	Long: `Usage: seed tables [[ ( in | except [in] ) ] <schema_name>[,<schema_name> ...]];

This command executes seed scripts in the following location:
- ./<schema_name>/tables/*.seed.*

If no <schema_name> is specified then all tables are seeded.

If the file extension is "sql" then the script is run directly on the database, otherwise
it is run as a shell script with "sh <seed_file>".

Examples:
  seed tables;
  seed tables in this_schema;
  seed tables except in that_schema;`,
	Run: runSeedCommand,
}

func init() {
	seedCmd.AddCommand(seedTablesCmd)
}
