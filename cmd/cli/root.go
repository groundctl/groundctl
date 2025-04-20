package main

import (
	"github.com/groundctl/groundctl/cmd/cli/stack"
	"github.com/groundctl/groundctl/cmd/cli/version"
	"github.com/groundctl/groundctl/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
)

var (
	verbose      int
	outputFormat output.Format
)

var RootCmd = &cobra.Command{
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
	// Add flags
	RootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "increase output verbosity")
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	RootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	RootCmd.PersistentFlags().VarP(
		enumflag.New(&outputFormat, "(normal|json)", output.FormatIds, enumflag.EnumCaseInsensitive),
		"format", "f",
		"the output format",
	)
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	// Add commands
	RootCmd.AddCommand(
		version.VersionCmd,
		stack.StackCmd,
	)
}
