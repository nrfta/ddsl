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

// viewsCmd represents the views command
var viewsCmd = &cobra.Command{
	Use:   "views",
	Short: "Create or drop views in a given schema",
	Long: `Usage: ddsl ( create | drop ) views [in] <schema_name> [ -W <exclude_view_name> ...]

Examples:
  ddsl create views in this_schema
  ddsl drop views that_schema -W not_that_view`,
	Run: createOrDropViews,
}

func createOrDropViews(cmd *cobra.Command, args []string) {
	fmt.Println("views called")
}

func init() {
	defineExcludeViewFlag(viewsCmd)
}
