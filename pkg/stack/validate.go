package stack

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

// Validate will ensure all fields in the stack template are valid
func (s *Stack) Validate() error {
	// Validate metadata fields
	if err := s.validateMetadata(); err != nil {
		return err
	}
	// Validate inputs
	if err := s.validateInputs(); err != nil {
		return err
	}
	// Validate layers (and steps)
	if err := s.validateLayers(); err != nil {
		return err
	}
	// Validate all templates
	if err := s.validateTemplates(); err != nil {
		return err
	}
	return nil
}

func (s *Stack) validateMetadata() error {
	logrus.WithField("stack", s.Name).Trace("Validating stack metadata")
	if s.Version == "" || s.Name == "" || s.Provider.Type == "" {
		return fmt.Errorf("stack is missing required metadata fields")
	}
	return nil
}

func (s *Stack) validateInputs() error {
	logrus.WithField("stack", s.Name).Trace("Validating stack inputs")
	for name, input := range s.Inputs {
		if input.Type == "" {
			return fmt.Errorf("input '%s' must have a type", name)
		}
		if !input.Required && input.Default == nil {
			return fmt.Errorf("input '%s' is not required but has no default value", name)
		}
		if len(input.Allowed) > 0 && input.Default != nil {
			allowed := false
			for _, val := range input.Allowed {
				if input.Default == val.Value {
					allowed = true
					break
				}
			}
			if !allowed {
				return fmt.Errorf("default value of input '%s' is not in the allowed list", name)
			}
		}
	}
	return nil
}

func (s *Stack) validateLayers() error {
	logrus.WithField("stack", s.Name).Trace("Validating stack layers")
	for _, layer := range s.Layers {
		if layer.Name == "" {
			return fmt.Errorf("layer is missing a name")
		}
		// Validate steps
		if err := layer.validateSteps(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Layer) validateSteps() error {
	logrus.WithField("layer", l.Name).Trace("Validating layer steps")
	for _, step := range l.Steps {
		if step.Name == "" || step.Action == "" {
			return fmt.Errorf("step in layer '%s' is missing name or action", l.Name)
		}
	}
	return nil
}

func (s *Stack) validateTemplates() error {
	logrus.WithField("stack", s.Name).Trace("Validating stack templates")
	var tmplCheck func(val any) error
	tmplCheck = func(val any) error {
		switch v := val.(type) {
		case string:
			// Check if string contains a template
			if !strings.Contains(v, "{{") {
				return nil
			}
			varPaths, err := ExtractVariablePaths(v)
			if err != nil {
				return err
			}

			// Validate each variable path
			for _, parts := range varPaths {
				if len(parts) == 0 {
					continue
				}
				root := parts[0]

				switch root {
				case "input":
					if len(parts) < 2 {
						return fmt.Errorf("invalid input reference: '%s' (missing key)", strings.Join(parts, "."))
					}
					if _, ok := s.Inputs[parts[1]]; !ok {
						return fmt.Errorf("undefined input '%s'", parts[1])
					}
				case "secret":
					if len(parts) < 2 {
						return fmt.Errorf("invalid secret reference: '%s' (missing key)", strings.Join(parts, "."))
					}
					if _, ok := s.Secrets[parts[1]]; !ok {
						return fmt.Errorf("undefined secret '%s'", parts[1])
					}
				default:
					// Treat everything else as a registered variable
					if _, ok := s.RegisteredVariables[root]; !ok {
						return fmt.Errorf("undefined registered variable: '%s'", root)
					}
				}
			}
			return nil
		case map[string]any:
			for _, val := range v {
				err := tmplCheck(val)
				if err != nil {
					return err
				}
			}
			return nil
		case []any:
			for i := range v {
				err := tmplCheck(v[i])
				if err != nil {
					return err
				}
			}
			return nil
		}
		return nil
	}

	// Validate provider properties
	err := tmplCheck(s.Provider.Properties)
	if err != nil {
		return fmt.Errorf("failed to resolve provider properties: %w", err)
	}

	// Validate each step’s params
	for i := range s.Layers {
		for j := range s.Layers[i].Steps {
			err := tmplCheck(s.Layers[i].Steps[j].Params)
			if err != nil {
				return fmt.Errorf("invalid template in step '%s' in layer '%s': %w", s.Layers[i].Steps[j].Name, s.Layers[i].Name, err)
			}
		}
	}

	return nil
}

func (s *Stack) validateTemplateSyntax(tmplStr string) error {
	// Validate the template syntax
	_, err := template.New("").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template syntax error: %w", err)
	}
	return nil
}

func (s *Stack) validateTemplateIdentifiers(tmplStr string) error {
	// Simple regex to extract expressions like {{ x.y.z | something }}
	re := regexp.MustCompile(`{{\s*([^{}\s|]+)`)
	matches := re.FindAllStringSubmatch(tmplStr, -1)

	for _, match := range matches {
		raw := match[1] // e.g., "input.foo", "my_var.bar"
		parts := strings.Split(raw, ".")
		if len(parts) == 0 {
			continue
		}

		// Check root identifier
		root := parts[0]
		if root == "input" || root == "secret" {
			if len(parts) < 2 {
				return fmt.Errorf("%s is referenced without key (e.g. '{{ %s.example_value }}')", root, root)
			}
			key := fmt.Sprintf("%s.%s", root, parts[1])
			if root == "input" {
				if _, ok := s.Inputs[key]; !ok {
					return fmt.Errorf("input '%s' not found", parts[1])
				}
			} else {
				if _, ok := s.Secrets[key]; !ok {
					return fmt.Errorf("secret '%s' not found", parts[1])
				}
			}
		} else {
			// Treat as registered variable
			if _, ok := s.RegisteredVariables[root]; !ok {
				return fmt.Errorf("references undefined registered variable '%s'", root)
			}
		}
	}

	return nil
}
