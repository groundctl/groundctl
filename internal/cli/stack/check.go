package stack

import (
	"errors"
	"fmt"
	"os"

	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CheckCmd struct{}

func (c *CheckCmd) Run(cmd *cobra.Command, args []string) error {
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
	logrus.Debug("Parsing stack...")
	parsedStack, err := stack.Parse(templateBytes)
	if err != nil {
		return fmt.Errorf("failed to parse stack template: %v", err)
	}
	// Validate the stack
	logrus.Debug("Validating stack...")
	if err = parsedStack.Validate(); err != nil {
		return fmt.Errorf("failed to validate stack: %v", err)
	}
	logrus.Infof("Stack file %q is valid!", args[0])
	return nil
}
