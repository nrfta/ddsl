package cmd

import (
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

// seedSchemaCmd represents the seedSchema command
var seedSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: parser.ShortDesc("seed schema"),
	Long: `Usage: seed schema <schema_name> [( with | without ) <seed_name>[,<seed_name> ...]];

This command executes seed scripts in the following location:
- ./<schema_name>/seeds/<seed_name>.*

If no <seed_name> is provided, then all seeds are executed.
If "with <seed_name>" is provided, then only the specified seed(s) are executed.
If "without <seed_name>" is provided, then all seeds are executed except the specified seed(s).

If the file extension is "sql" then the script is run directly on the database, otherwise
it is run as a shell script with "sh <seed_file>".

Examples:
  seed schema this_schema;
  seed schema this_schema with this_seed;
  seed schema that_schema without that_seed;`,
	Run: runSeedCommand,
}

func init() {
	seedCmd.AddCommand(seedSchemaCmd)
}
