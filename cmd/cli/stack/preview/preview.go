package preview

import (
	"github.com/groundctl/groundctl/internal/cli/stack"
	"github.com/spf13/cobra"
)

var c stack.PreviewCmd

var PreviewCmd = &cobra.Command{
	Use:   "preview filename",
	Short: "Preview a stack.",
	Long: `Preview a stack.

Stacks are groundctl's environment templates.`,
	Aliases: []string{"p"},
	Args:    cobra.ExactArgs(1),
	Example: "groundctl preview example.stack",
	RunE:    c.Run,
}
