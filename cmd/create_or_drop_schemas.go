package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// schemasCmd represents the schemas command
var schemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: "Create or drop all schemas in the database",
	Long: `Usage: ( create | drop ) schemas [ -N <exclude_schema_name> ... ]

Examples:
  create schemas
  drop schemas -N not_this_schema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("schemas called")
	},
}

func init() {
	defineExcludeSchemaFlag(schemasCmd)
}

