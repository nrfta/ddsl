package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// constraintsCmd represents the constraints command
var constraintsCmd = &cobra.Command{
	Use:   "constraints",
	Short: "Create or drop constraints on a given table",
	Long: `Usage: ( create | drop ) constraints [on] <table_name>

Examples:
  ddsl create constraints on this_schema.this_table
  ddsl drop constraints that_schema.that_table`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("constraints called")
	},
}
