package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// viewsCmd represents the views command
var viewsCmd = &cobra.Command{
	Use:   "views",
	Short: "Create or drop views in a given schema",
	Long: `Usage: ( create | drop ) views [in] <schema_name> [ -W <exclude_view_name> ...]

Examples:
  create views in this_schema
  drop views that_schema -W not_that_view`,
	Run: createOrDropViews,
}

func createOrDropViews(cmd *cobra.Command, args []string) {
	fmt.Println("views called")
}

func init() {
	defineExcludeViewFlag(viewsCmd)
}
