package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// indexesCmd represents the indexes command
var indexesCmd = &cobra.Command{
	Use:   "indexes",
	Short: "Create or drop indexes on a given table or view",
	Long: `Usage: ( create | drop ) indexes [on] <table_or_view_name>

Examples:
  create indexes on this_schema.this_table
  drop indexes that_schema.that_view`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("indexes called")
	},
}
