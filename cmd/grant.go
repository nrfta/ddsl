package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/parser"
	"os"

	"github.com/spf13/cobra"
)

// grantCmd represents the grant command
var grantCmd = &cobra.Command{
	Use:   "grant",
}

var grantPrivilegesCmd = &cobra.Command{
	Use: "privileges",
	Aliases: []string{"privs"},
}

var grantPrivilegesOnCmd = &cobra.Command{
	Use: "on",
	Short: parser.ShortDesc("grant privileges on"),
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)
	grantPrivilegesCmd.AddCommand(grantPrivilegesOnCmd)
	grantPrivilegesOnCmd.AddCommand(grantDatabaseCmd)
	grantPrivilegesOnCmd.AddCommand(grantSchemaCmd)
	grantPrivilegesOnCmd.AddCommand(grantSchemasCmd)
	grantPrivilegesOnCmd.AddCommand(grantTableCmd)
	grantPrivilegesOnCmd.AddCommand(grantTablesCmd)
	grantPrivilegesOnCmd.AddCommand(grantViewCmd)
	grantPrivilegesOnCmd.AddCommand(grantViewsCmd)
}
