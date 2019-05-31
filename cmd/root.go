package cmd

import (
	"fmt"
	"github.com/neighborly/ddsl/repl"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Set this value using the Go linker based on git tag. See https://stackoverflow.com/a/11355611
var appVersion string

var (
	version bool
	command string
	file string
	schemas []string
	tables []string
	views []string
	excludeSchemas []string
	excludeTables []string
	excludeViews []string
)
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ddsl [OPTIONS] [COMMAND]",
	Short: `Data-Definition-Specific Language (DDSL, pronounced "diesel") 
provides a scripting language for DDL and migrations.`,
	Long: `ddsl executes commands written in DDSL. Commands can either be
one-off or stored in a ddsl file. In addition, ddsl files may made directly
executable.

Run the REPL shell:
	ddsl

Run ddsl commands from the command line:
    ddsl [OPTIONS] command

Run commands from a ddsl file:
    ddsl [OPTIONS] -f /path/to/file.ddsl

Make ddsl file executable with "chmod +x file.ddsl" and adding shebang.
Requires environment variables to set options. The "ddsl" command is
omitted from the beginning of each line.
    #!/usr/bin/env ddsl
	begin transaction;
    command; command;
	command;
	commit transaction;
    etc...`,

	Run: func(cmd *cobra.Command, args []string) {
		db := viper.GetString("database")
		src := viper.GetString("source")
		exitCode := 0
		switch {

		case version:
			fmt.Println("ddsl version", appVersion)

		case len(file) > 0:
			ensureArgs(src, db)
			ec, err := runFile(src, db, file)
			if err != nil {
				fmt.Println(err)
			}
			exitCode = ec

		default:
			ec, err := repl.Run(src, db)
			if err != nil {
				fmt.Println(err)
			}
			exitCode = ec
		}

		os.Exit(exitCode)
	},
}

func ensureArgs(src string, db string) {
	if len(db) == 0 {
		fmt.Println("no database URL provided")
		os.Exit(1)
	}
	if len(src) == 0 {
		fmt.Println("no source repository provided")
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetEnvPrefix("ddsl")
	viper.BindEnv("source")
	viper.BindEnv("database")

	rootCmd.PersistentFlags().StringP("source", "s", "", "DDL source repo (default DDSL_SOURCE)")
	viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))

	rootCmd.PersistentFlags().StringP("database", "d", "", "URL for RDS and database (default DDSL_DATABASE)")
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))

	rootCmd.PersistentFlags().Bool("dry-run", false, "take no action but output what would be done")
	viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))

	rootCmd.Flags().BoolVar(&version, "version", false, "show version number and exit")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "file containing DDSL commands")

}

func addCreateOrDropSubCmds(createOrDropCmd *cobra.Command) {
	createOrDropCmd.AddCommand(databaseCmd)
	createOrDropCmd.AddCommand(rolesCmd)
	createOrDropCmd.AddCommand(extensionsCmd)
	createOrDropCmd.AddCommand(foreignKeysCmd)
	createOrDropCmd.AddCommand(schemasCmd)
	createOrDropCmd.AddCommand(typesCmd)
	createOrDropCmd.AddCommand(schemaCmd)
	createOrDropCmd.AddCommand(tablesCmd)
	createOrDropCmd.AddCommand(viewsCmd)
	createOrDropCmd.AddCommand(tableCmd)
	createOrDropCmd.AddCommand(viewCmd)
	createOrDropCmd.AddCommand(constraintsCmd)
	createOrDropCmd.AddCommand(indexesCmd)
}


func addGrantOrRevokeSubCmds(grantOrRevokeCmd *cobra.Command) {
	grantOrRevokeCmd.AddCommand(grantPrivilegesCmd)
	grantOrRevokeCmd.AddCommand(grantPrivilegesOnCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
}
