/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adventofcode2025",
	Short: "A brief description of your application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		extraVerbose, _ := cmd.Flags().GetBool("extra-verbose")

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		if extraVerbose {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Logger()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("follow-up", "f", false, "runs the follow up")
	rootCmd.PersistentFlags().StringP("input-file", "i", "inputs/01", "select file to parse")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolP("extra-verbose", "t", false, "enable trace logging")
}
