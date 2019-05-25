package cmd

import (
	"fmt"
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

Run ddsl commands:
    ddsl [OPTIONS] COMMAND1

Run commands from a ddsl file:
    ddsl [OPTIONS] -f /path/to/file.ddsl

Make ddsl file executable with "chmod +x file.ddsl" and adding shebang.
Requires environment variables to set options. The "ddsl" command is
omitted from the beginning of each line.
    #!/usr/bin/env ddsl
    COMMAND
	COMMAND
    etc...`,

	Run: func(cmd *cobra.Command, args []string) {
		db := viper.GetString("database")
		src := viper.GetString("source")

		switch {

		case version:
			fmt.Println("ddsl version", appVersion)
			os.Exit(0)

		case len(file) > 0:
			ensureArgs(src, db)
			exitCode, err := runFile(src, db, file)
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(exitCode)

		default:
			fmt.Println("[ERROR] unknown usage")
			cmd.Usage()
			os.Exit(1)
		}
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

func defineTableFlag(command *cobra.Command) {
	command.PersistentFlags().StringSliceVarP(&tables,"table", "t", nil, "table to operate upon with schema optionally specified ([schema.]table). Can be specified more than once.")
}
func defineExcludeTableFlag(command *cobra.Command) {
	command.PersistentFlags().StringSliceVarP(&excludeTables,"exclude-table", "T", nil, "table to exclude with schema optionally specified ([schema.]table). Can be specified more than once.")
}

func defineViewFlag(command *cobra.Command) {
	command.PersistentFlags().StringSliceVarP(&views, "view", "w", nil, "view to operate upon with schema optionally specified ([schema.]view). Can be specified more than once.")
}
func defineExcludeViewFlag(command *cobra.Command) {
	command.PersistentFlags().StringSliceVarP(&excludeViews,"exclude-view", "W", nil, "view to exclude with schema optionally specified ([schema.]view). Can be specified more than once.")
}

func defineSchemaFlag(command *cobra.Command) {
	command.PersistentFlags().StringSliceVarP(&schemas,"schema", "n", nil, "schema to operate upon")
}
func defineExcludeSchemaFlag(command *cobra.Command) {
	command.PersistentFlags().StringSliceVarP(&excludeSchemas,"exclude-schema", "N", nil, "schema to exclude")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
}
