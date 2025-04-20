package stack

import (
	"github.com/groundctl/groundctl/cmd/cli/stack/check"
	"github.com/groundctl/groundctl/cmd/cli/stack/preview"
	"github.com/spf13/cobra"
)

var StackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Work with stacks",
	Long: `Work with stacks. Perform operations like validating templates, deploying stacks, etc.

Stacks are groundctl's environment templates.`,
	Aliases: []string{"s", "stacks"},
}

func init() {
	StackCmd.AddCommand(
		check.CheckCmd,
		preview.PreviewCmd,
	)
}
