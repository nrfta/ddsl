package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seedDatabaseCmd represents the seedDatabase command
var seedDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Seeds the database with all seeds",
	Long: `Usage: ddsl seed database [-n <include_schema_name> ...] [-N <exclude_schema_name> ...] \
                          [-t <include_table_name> ...] [-T <exclude_table_name> ...]

This command traverses the repo and executes seed scripts in the following locations:
- ./seeds/*
- ./schemas/seeds/*
- ./schemas/tables/*.seed.*

Examples:
  seed database
  seed database -n this_schema -n that_schema -t other_schema.other_table
  seed database -N not_this_schema -T that_schema.not_that_table`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seedDatabase called")
	},
}

func init() {
	seedCmd.AddCommand(seedDatabaseCmd)
	defineSchemaFlag(seedDatabaseCmd)
	defineTableFlag(seedDatabaseCmd)
	defineExcludeSchemaFlag(seedDatabaseCmd)
	defineExcludeTableFlag(seedDatabaseCmd)
}
