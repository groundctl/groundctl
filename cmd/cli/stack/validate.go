package stack

import (
	"github.com/groundctl/groundctl/internal/cli/stack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var c stack.ValidateCmd

var validateCmd = &cobra.Command{
	Use:   "validate filename",
	Short: "Parse and validate a stack template file.",
	Long: `Parse and validate a stack template file.

Stacks are groundctl's environment templates.`,
	Aliases: []string{"v"},
	Args:    cobra.ExactArgs(1),
	Example: "groundctl validate example.stack",
	RunE:    c.Run,
}

func init() {
	validateCmd.Flags().BoolVar(&c.SkipDependencies, "skip-references", false, "Skip validating all template references (e.g. \"${var.my_var}\")")
	viper.BindPFlag("skip-references", validateCmd.Flags().Lookup("skip-references"))
	validateCmd.Flags().BoolVar(&c.SkipDependencies, "skip-dependencies", false, "Skip validating all template dependencies (i.e. depends_on lists)")
	viper.BindPFlag("skip-dependencies", validateCmd.Flags().Lookup("skip-dependencies"))
}
