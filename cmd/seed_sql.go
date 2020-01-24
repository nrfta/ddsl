package cmd

import (
	"github.com/nrfta/ddsl/parser"

	"github.com/spf13/cobra"
)

// seedSqlCmd represents the seedSql command
var seedSqlCmd = &cobra.Command{
	Use:   "sql",
	Short: parser.ShortDesc("seed sql"),
	Long: `The seed sql command is functionally identical to the sql command.

Examples:
  seed sql -f ./seeds/seed1.sql -f ./seeds/seed2.sql;
  seed sql ` + "`SQL script as text`;" + `

Note that SQL script is enclosed in backticks and can be multiple lines`,
	Run: runSeedCommand,
}

func init() {
	seedCmd.AddCommand(seedSqlCmd)

	seedSqlCmd.PersistentFlags().StringSliceVarP(&sqlFiles, "file","f", nil, "SQL file path, may be provided more than once")
}
