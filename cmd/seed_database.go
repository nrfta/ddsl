package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

// seedDatabaseCmd represents the seedDatabase command
var seedDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: parser.ShortDesc("seed database"),
	Long: `Usage: seed database [( with | without ) <seed_name>[,<seed_name> ...]];

This command executes seed scripts in the following location:
- ./seeds/<seed_name>.*

If the file extension is "sql" then the script is run directly on the database, otherwise
it is run as a shell script with "sh <seed_file>".

Examples:
  seed database;
  seed database with this_seed;
  seed database without that_seed;`,
	Run: runSeedCommand,
}

func init() {
	seedCmd.AddCommand(seedDatabaseCmd)
}
