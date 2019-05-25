package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seedTableCmd represents the seedTable command
var seedTableCmd = &cobra.Command{
	Use:   "table",
	Short: "Seeds the specified table",
	Long: `Usage: ddsl seed table <table_name>

Examples:
  seed table this_schema.this_table`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seedTable called")
	},
}

func init() {
	seedCmd.AddCommand(seedTableCmd)
}
