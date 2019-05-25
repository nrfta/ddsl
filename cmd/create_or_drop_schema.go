package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Create or drop a given schema",
	Long: `Examples:
  ddsl create schema this_schema
  ddsl drop schema that_schema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("schema called")
	},
}
