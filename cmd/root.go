package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Set this value using the Go linker based on git tag. See https://stackoverflow.com/a/11355611
var appVersion string

var version bool
var command string
var file string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ddsl [OPTIONS] [COMMAND]",
	Short: `Data-Definition-Specific Language (DDSL, pronounced "diesel") provides a scripting language for DDL and migrations.`,
	Long: `ddsl executes commands written in DDSL. Commands can either be one-off or stored in a ddsl file. 
In addition, ddsl files may made directly executable.

Run a one-off command:
    ddsl [OPTIONS] -c "COMMAND"

Run commands from a ddsl file:
    ddsl [OPTIONS] -f /path/to/file.ddsl

Make ddsl file executable with "chmod +x file.ddsl" and adding shebang:
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

		case len(command) > 0 && len(file) > 0:
			fmt.Println("[ERROR] file and command arguments cannot be used together")
			os.Exit(1)

		case len(command) > 0:
			ensureArgs(src, db)
			exitCode, err := runCommand(src, db, command)
			if err != nil {
				fmt.Printf("[ERROR] %s", err)
			}
			os.Exit(exitCode)

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

	rootCmd.PersistentFlags().BoolVar(&version, "version", false, "show version number and exit")
	rootCmd.PersistentFlags().StringVarP(&command, "command", "c", "", "DDSL command to run")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "file containing DDSL commands")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
}
