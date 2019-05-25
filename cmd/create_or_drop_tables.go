package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tablesCmd represents the tables command
var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Create or drop tables in a given schema",
	Long: `Usage: ddsl ( create | drop ) tables [in] <schema_name> [ -T <exclude_table_name> ...]

Examples:
  ddsl create tables in this_schema
  ddsl drop tables that_schema -T not_that_table`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tables called")
	},
}

func init() {
	defineExcludeTableFlag(tablesCmd)
}
