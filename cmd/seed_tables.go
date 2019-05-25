package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// seedTablesCmd represents the seedTables command
var seedTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Seeds all tables in the given schema",
	Long: `Usage: ddsl seed tables [in] <schema_name> [-T exclude_table_name ...]

Examples:
  seed tables in this_schema
  seed tables that_schema -T except_that_table`,
	Run: seedTables,
}

func seedTables(cmd *cobra.Command, args []string) {
	fmt.Println("seed tables called")
}

func init() {
	seedCmd.AddCommand(seedTablesCmd)

	defineExcludeTableFlag(seedTablesCmd)
}
