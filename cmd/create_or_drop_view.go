package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Create or drop a given view",
	Long: `Examples:
  ddsl create view this_schema.this_view
  ddsl drop view that_schema.that_view`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("view called")
	},
}
