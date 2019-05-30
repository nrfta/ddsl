package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"

	"github.com/spf13/cobra"
)

var sqlFiles []string

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: parser.ShortDesc("sql"),
	Long: `The sql command is functionally identical to the seed sql command.
Both are provided to enable more intention-revealing scripts.

Examples:
  sql -f /path/to/sql_file1.sql -f /path/to/sql_file2.sql;
  sql ` + "`SQL script as text`;" + `

Note that SQL script is enclosed in backticks and can be multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sql called")
	},
}

func init() {
	rootCmd.AddCommand(sqlCmd)

	sqlCmd.PersistentFlags().StringSliceVarP(&sqlFiles, "file","f", nil, "SQL file path, may be provided more than once")

}
