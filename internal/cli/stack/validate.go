package stack

import (
	"errors"
	"fmt"
	"os"

	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ValidateCmd struct {
	SkipDependencies bool
	SkipReferences   bool
}

func (c *ValidateCmd) Run(cmd *cobra.Command, args []string) error {
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
	if c.SkipDependencies {
		logrus.Warn("Skipping dependency checks...")
	} else {
		logrus.Debug("Validating dependencies...")
		if err = parsedStack.ValidateDependsOn(); err != nil {
			return fmt.Errorf("failed to validate stack resource dependencies: %v", err)
		}
	}
	// Validate the stack variable references
	if c.SkipReferences {
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
