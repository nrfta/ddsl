package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seedSchemaCmd represents the seedSchema command
var seedSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Seeds the specified schema",
	Long: `Examples:
  seed schema this_schema`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seedSchema called")
	},
}

func init() {
	seedCmd.AddCommand(seedSchemaCmd)
}
