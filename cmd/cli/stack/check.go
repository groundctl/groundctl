package stack

import (
	"github.com/groundctl/groundctl/internal/cli/stack"
	"github.com/spf13/cobra"
)

var c stack.CheckCmd

var checkCmd = &cobra.Command{
	Use:   "check filename",
	Short: "Parse and validate a stack template file.",
	Long: `Parse and validate a stack template file.

Stacks are groundctl's environment templates.`,
	Aliases: []string{"c"},
	Args:    cobra.ExactArgs(1),
	Example: "groundctl check example.stack",
	RunE:    c.Run,
}

func init() {
}
