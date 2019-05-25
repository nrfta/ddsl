package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// extensionsCmd represents the extensions command
var extensionsCmd = &cobra.Command{
	Use:   "extensions",
	Short: "Create or drop extensions",
	Long: `Usage: ddsl ( create | drop ) extensions [in] schema_name

Examples:
  ddsl create extensions in this_schema
  ddsl drop extensions that_schema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("extensions called")
	},
}
