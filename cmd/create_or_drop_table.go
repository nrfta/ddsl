package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tableCmd represents the table command
var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Create or drop a given table",
	Long: `Examples:
  ddsl create table this_schema.this_table
  ddsl drop table that_schema.that_table`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("table called")
	},
}
