package cmd

import (
	"fmt"
	"github.com/nrfta/ddsl/log"
	"github.com/nrfta/ddsl/parser"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	GRANT = "grant"
	REVOKE = "revoke"
)

// grantCmd represents the grant command
var grantCmd = &cobra.Command{
	Use:   "grant",
	Short: parser.ShortDesc("grant"),
	Long: `Grant or revoke privileges on database objects`,
}

var grantPrivilegesCmd = &cobra.Command{
	Use: "privileges",
	Short: parser.ShortDesc("grant privileges"),
	Long: `Grant or revoke privileges on database objects`,
	Aliases: []string{"privs"},
}

var grantPrivilegesOnCmd = &cobra.Command{
	Use: "on",
	Short: parser.ShortDesc("grant privileges on"),
	Long: `Grant or revoke privileges on database objects`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("additional arguments required, use -h for help")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)
	grantCmd.AddCommand(grantPrivilegesCmd)
	grantPrivilegesCmd.AddCommand(grantPrivilegesOnCmd)
	grantPrivilegesOnCmd.AddCommand(grantDatabaseCmd)
	grantPrivilegesOnCmd.AddCommand(grantSchemaCmd)
	grantPrivilegesOnCmd.AddCommand(grantSchemasCmd)
	grantPrivilegesOnCmd.AddCommand(grantTableCmd)
	grantPrivilegesOnCmd.AddCommand(grantTablesCmd)
	grantPrivilegesOnCmd.AddCommand(grantViewCmd)
	grantPrivilegesOnCmd.AddCommand(grantViewsCmd)
	grantPrivilegesOnCmd.AddCommand(grantFunctionCmd)
	grantPrivilegesOnCmd.AddCommand(grantFunctionsCmd)
	grantPrivilegesOnCmd.AddCommand(grantProcedureCmd)
	grantPrivilegesOnCmd.AddCommand(grantProceduresCmd)
}

func grantOrRevoke(cmd *cobra.Command) string {
	for _, a := range os.Args {
		switch a {
		case GRANT:
			return GRANT
		case REVOKE:
			return REVOKE
		}
	}
	return ""
}

func runGrantOrRevokeCommand(cmd *cobra.Command, args []string) {
	gOrR := grantOrRevoke(cmd)
	if len(gOrR) == 0 {
		panic("use only for grant or revoke")
	}
	command := fmt.Sprintf("%s privileges on %s", gOrR, cmd.Use)
	if len(args) > 0 {
		command += " "
	}
	command += strings.Join(args, " ")

	code, err := runCLICommand(command)
	if err != nil {
		log.Error(err.Error())
	}
	os.Exit(code)
}
