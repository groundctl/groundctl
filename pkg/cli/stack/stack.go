package stack

import (
	"context"

	"github.com/bep/simplecobra"
)

func NewStackCommand() *stackCommand {
	return &stackCommand{}
}

type stackCommand struct{}

func (c *stackCommand) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{
		newValidateCommand(),
	}
}

func (c *stackCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	return nil
}

func (c *stackCommand) Name() string {
	return "stacks"
}

func (c *stackCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	cd.CobraCommand.Help()
	return nil
}

func (c *stackCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Work with stacks"
	cmd.Long = `Work with stacks. Perform operations like validating templates, deploying stacks, etc.

Stacks are groundctl's environment templates.`
	cmd.Aliases = []string{"s", "stacks"}
	return nil
}
