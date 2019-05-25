// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seedSqlCmd represents the seedSql command
var seedSqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Seeds the database through a SQL command",
	Long: `The seed sql command is functionally identical to the sql command.

Examples:
  seed sql -f ./seeds/seed1.sql -f ./seeds/seed2.sql
  seed sql ` + "`SQL script as text`" + `

Note that SQL script is enclosed in backticks and can be multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seedSql called")
	},
}

func init() {
	seedCmd.AddCommand(seedSqlCmd)

	seedSqlCmd.PersistentFlags().StringSliceVarP(&sqlFiles, "file","f", nil, "SQL file path, may be provided more than once")
}
