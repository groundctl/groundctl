package stack

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type PreviewCmd struct{}

func (c *PreviewCmd) Run(cmd *cobra.Command, args []string) error {
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
	outputPreview(parsedStack)
	return nil
}

func outputPreview(s *stack.Stack) {
	fmt.Printf("Stack Plan Preview: %s (version: %s)\n", s.DisplayName, s.Version)
	if s.Description != "" {
		fmt.Printf("Description: %s\n", s.Description)
	}
	fmt.Println(strings.Repeat("=", 50))

	// Provider
	fmt.Printf("\n🔌 Provider: %s\n", s.Provider.Type)
	if len(s.Provider.Properties) > 0 {
		keys := sortedKeys(s.Provider.Properties)
		for _, k := range keys {
			fmt.Printf("  - %s: %s\n", k, s.Provider.Properties[k])
		}
	}

	// Inputs
	if len(s.Inputs) > 0 {
		fmt.Println("\n📥 Inputs:")
		keys := sortedKeys(s.Inputs)
		for _, name := range keys {
			input := s.Inputs[name]
			fmt.Printf("  - %s (%s)", name, input.Type)
			if input.Required {
				fmt.Print(" [required]")
			}
			if input.Default != nil {
				fmt.Printf(" [default: %v]", input.Default)
			}
			fmt.Println()
			if input.Label != "" {
				fmt.Printf("      → Label: %s\n", input.Label)
			}
			if input.Description != "" {
				fmt.Printf("      → Description: %s\n", input.Description)
			}
			if len(input.Allowed) > 0 {
				fmt.Print("      → Allowed:\n")
				for _, v := range input.Allowed {
					fmt.Printf("          - %v (%s)\n", v.Value, v.Label)
				}
			}
		}
	}

	// Secrets
	if len(s.Secrets) > 0 {
		fmt.Println("\n🔐 Secrets:")
		keys := sortedKeys(s.Secrets)
		for _, name := range keys {
			secret := s.Secrets[name]
			fmt.Printf("  - %s (%s)\n", name, secret.Type)
			if secret.Label != "" {
				fmt.Printf("      → Label: %s\n", secret.Label)
			}
			if secret.Description != "" {
				fmt.Printf("      → Description: %s\n", secret.Description)
			}
			if len(secret.Allowed) > 0 {
				fmt.Print("      → Allowed:\n")
				for _, v := range secret.Allowed {
					fmt.Printf("          - %v (%s)\n", v.Value, v.Label)
				}
			}
		}
	}

	// Layers & Steps
	if len(s.Layers) > 0 {
		fmt.Println("\n🧱 Layers:")
		for i, layer := range s.Layers {
			fmt.Printf("  %d. %s\n", i+1, layer.Name)
			for j, step := range layer.Steps {
				action := step.Action
				if action == "" {
					// fallback to raw action key
					for k := range step.Raw {
						action = k
						break
					}
				}
				fmt.Printf("     %d.%d Step: %s [%s]\n", i+1, j+1, step.Name, action)

				if step.Register != "" {
					fmt.Printf("       → registers: %s\n", step.Register)
				}
				if len(step.Tags) > 0 {
					fmt.Printf("       Tags: %v\n", step.Tags)
				}
			}
		}
	}

	// Outputs
	if len(s.Outputs) > 0 {
		fmt.Println("\n📤 Outputs:")
		keys := sortedKeys(s.Outputs)
		for _, name := range keys {
			out := s.Outputs[name]
			fmt.Printf("  - %s = %s\n", name, out.Value)
			if out.Description != "" {
				fmt.Printf("      → %s\n", out.Description)
			}
		}
	}
}

func sortedKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})
	return keys
}
