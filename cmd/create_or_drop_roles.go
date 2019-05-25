package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Create or drop roles",
	Long: `Use caution: roles may be shared across databases.

Examples:
  create roles
  drop roles`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("roles called")
	},
}
