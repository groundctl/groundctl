package stack

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bep/simplecobra"
	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ()

func newValidateCommand() *validateCommand {
	return &validateCommand{}
}

type validateCommand struct {
	skipDependencies bool
	skipReferences   bool
}

func (c *validateCommand) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{}
}

func (c *validateCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	return nil
}

func (c *validateCommand) Name() string {
	return "validate"
}

func (c *validateCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("accepts 1 arg(s), received %d", len(args))

	}
	// Read the template file
	templateBytes, err := os.ReadFile(args[0])
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file does not exist")
		}
		return fmt.Errorf("failed to read file: %v", err)
	}
	// Parse as a template
	parsedStack, err := stack.Parse(templateBytes)
	if err != nil {
		return fmt.Errorf("failed to parse stack template: %v", err)
	}
	// Validate the stack dependencies
	if c.skipDependencies {
		logrus.Warn("Skipping dependency checks...")
	} else {
		logrus.Debug("Validating dependencies...")
		if err = parsedStack.ValidateDependsOn(); err != nil {
			return fmt.Errorf("failed to validate stack resource dependencies: %v", err)
		}
	}
	// Validate the stack variable references
	if c.skipReferences {
		logrus.Warn("Skipping reference checks...")
	} else {
		logrus.Debug("Validating references...")
		if err = parsedStack.ValidateReferences(); err != nil {
			return fmt.Errorf("failed to validate stack references: %v", err)
		}
	}
	logrus.Infof("Stack file %q is valid!", args[0])
	return nil
}

func (c *validateCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Parse and validate a stack template file."
	cmd.Long = `Parse and validate a stack template file.

Stacks are groundctl's environment templates.`
	cmd.Aliases = []string{"v"}
	cmd.Example = "groundctl validate example.stack"
	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().BoolVar(&c.skipDependencies, "skip-references", false, "Skip validating all template references (e.g. \"${var.my_var}\")")
	cmd.Flags().BoolVar(&c.skipDependencies, "skip-dependencies", false, "Skip validating all template dependencies (i.e. depends_on lists)")
	return nil
}
