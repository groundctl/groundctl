package stack

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Stack, error) {
	var stack Stack
	err := yaml.Unmarshal(data, &stack)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stack template: %w", err)
	}

	// Validate stack resources
	if err := stack.ValidateReferences(); err != nil {
		return nil, fmt.Errorf("stack validation failed: %w", err)
	}

	return &stack, nil
}
