package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appVersion = ""

var source string
var database string
var version bool
var command string
var file string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ddsl [OPTIONS] [COMMAND]",
	Short: `Data-Definition-Specific Language (DDSL, pronounced "diesel") provides a scripting language for DDL and migrations.`,
	Long: `ddsl executes commands written in DDSL. Commands can either be one-off or stored in a ddsl file. In addition,
ddsl files may made directly executable.

Run a one-off command:
    ddsl [OPTIONS] -c "COMMAND"

Run commands from a ddsl file:
    ddsl [OPTIONS] -f /path/to/file.ddsl`,

	Run: func(cmd *cobra.Command, args []string) {
		switch {

		case version:
			fmt.Println("ddsl version", appVersion)
			os.Exit(0)

		case len(command) > 0 && len(file) > 0:
			fmt.Println("File and command arguments cannot be used together")
			os.Exit(1)

		case len(command) > 0:
			ensureArgs()
			exitCode, err := runCommand(source, database, command)
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(exitCode)

		case len(file) > 0:
			ensureArgs()
			exitCode, err := runFile(source, database, file)
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(exitCode)

		default:
			fmt.Println("Unknown usage")
			cmd.Usage()
			os.Exit(1)
		}
	},
}

func ensureArgs() {
	if len(database) == 0 {
		fmt.Println("no database URL provided")
		os.Exit(1)
	}
	if len(source) == 0 {
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

	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "DDL source repo (default DDSL_SOURCE)")
	viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))

	rootCmd.PersistentFlags().StringVarP(&source, "database", "d", "", "URL for RDS and database (default DDSL_DATABASE)")
	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))

	rootCmd.PersistentFlags().BoolVar(&version, "version", false, "show version number and exit")
	rootCmd.PersistentFlags().StringVarP(&command, "command", "c", "", "DDSL command to run")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "file containing DDSL commands")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("ddsl")
	viper.BindEnv("source")
	viper.BindEnv("database")
	viper.AutomaticEnv()

}
