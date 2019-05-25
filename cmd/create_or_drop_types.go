package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// typesCmd represents the types command
var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "Create or drop types in a given schema",
	Long: `Usage: ( create | drop ) types [in] <schema_name>

Examples:
  create types in this_schema
  drop types that_schema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("types called")
	},
}
