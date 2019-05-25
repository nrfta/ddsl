package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Creates or drops the database",
	Long: `This command creates the database itself. It does not create any of the
objects within the database.

Examples:
  create database
  drop database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("database called")
	},
}
