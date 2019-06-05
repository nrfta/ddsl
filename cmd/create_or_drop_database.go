package cmd

import (
	"github.com/neighborly/ddsl/parser"
	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: parser.ShortDesc("create database"),
	Long: `This command creates the database itself. It does not create any of the
objects within the database.

Examples:
  create database;
  drop database;`,
	Run: runCreateOrDropCommand,
}
