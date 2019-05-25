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

var upNumVersions uint

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrates the database up a given number of versions",
	Long: `Examples:
  migrate up -n 2`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("up called")
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)

	upCmd.PersistentFlags().UintVarP(&upNumVersions, "num-versions", "n", 0, "number of versions to migrate up")
}
