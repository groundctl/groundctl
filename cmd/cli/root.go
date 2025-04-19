package main

import (
	"github.com/groundctl/groundctl/cmd/cli/stack"
	"github.com/groundctl/groundctl/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
)

var (
	verbose      int
	outputFormat output.Format
)

var rootCmd = &cobra.Command{
	Use:   "groundctl",
	Short: "groundctl (Ground Control) is the CLI client for the groundctl server",
	Long:  "groundctl (Ground Control) is the CLI client for the groundctl server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Setup logging
		output.SetFormat(outputFormat)
		output.SetVerbosity(verbose)
	},
}

func init() {
	rootCmd.AddCommand(
		stack.StackCmd,
		versionCmd,
	)
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "increase output verbosity")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&outputFormat, "(normal|json)", output.FormatIds, enumflag.EnumCaseInsensitive),
		"format", "f",
		"the output format",
	)
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
