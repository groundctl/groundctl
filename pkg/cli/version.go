package cli

import (
	"context"
	"fmt"

	"github.com/bep/simplecobra"
	"github.com/groundctl/groundctl/internal/version"
)

func newVersionCommand() *versionCommand {
	return &versionCommand{}
}

type versionCommand struct{}

func (c *versionCommand) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{}
}

func (c *versionCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	return nil
}

func (c *versionCommand) Name() string {
	return "version"
}

func (c *versionCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	fmt.Printf("groundctl v%s\n", version.Version)
	return nil
}

func (c *versionCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Display version"
	cmd.Long = "Display version and environment info."
	return nil
}
