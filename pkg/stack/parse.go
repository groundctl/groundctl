package stack

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	// All of the reserved variables that cannot be used as registered variable names
	reservedVariables = []string{"input", "secret"}
)

// Parse reads in a stack template file and parses into a Stack struct
func Parse(data []byte) (*Stack, error) {
	var stack Stack
	if err := yaml.Unmarshal(data, &stack); err != nil {
		return nil, err
	}

	// Parse the actions out of all the steps
	for i := range stack.Layers {
		for j := range stack.Layers[i].Steps {
			// Parse the action names from the steps
			if err := stack.Layers[i].Steps[j].parseAction(); err != nil {
				return nil, err
			}
			// Register the variable name as a placeholder in the stack
			if stack.Layers[i].Steps[j].Register != "" {
				varName := stack.Layers[i].Steps[j].Register
				// Check variable name isn't a reserved name
				for _, reservedVarName := range reservedVariables {
					if varName == reservedVarName {
						return nil, fmt.Errorf("'%s' is a reserved variable name", varName)
					}
				}
				// Check if this variable has already been defined elsewhere
				if _, ok := stack.RegisteredVariables[varName]; ok {
					return nil, fmt.Errorf("variable '%s' in step '%s' already defined", varName, stack.Layers[i].Steps[j].Name)
				}
				stack.RegisteredVariables[varName] = make(map[string]any)
			}
		}
	}

	return &stack, nil
}

// Parses the action from the step definition
//
//	// "aws.vpc" is the action
//	name: Create VPC
//	aws.vpc:
//	  name: my-vpc
//	  cidr_block: 10.0.0.0/16
//	register: my_vpc
func (t *Step) parseAction() error {
	for k, v := range t.Raw {
		switch k {
		// Ignore all reserved attributes
		case "name", "register", "tags":
			continue
		// If not a reserved attribute, assume it's the action and parse
		default:
			t.Action = k
			params, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("invalid format for step action '%s' in step '%s'", k, t.Name)
			}
			t.Params = params
			return nil
		}
	}
	return fmt.Errorf("no action found in step '%s'", t.Name)
}
