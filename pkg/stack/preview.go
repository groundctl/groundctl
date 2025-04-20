package stack

import (
	"fmt"
	"sort"
	"strings"
)

func (s *Stack) PrintPreview() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Stack Plan Preview: %s (version: %s)\n", s.DisplayName, s.Version)
	if s.Description != "" {
		fmt.Fprintf(&b, "Description: %s\n", s.Description)
	}
	fmt.Fprintln(&b, strings.Repeat("=", 50))

	// Provider
	fmt.Fprintf(&b, "\n🔌 Provider: %s\n", s.Provider.Type)
	if len(s.Provider.Properties) > 0 {
		keys := sortedKeys(s.Provider.Properties)
		for _, k := range keys {
			fmt.Fprintf(&b, "  - %s: %v\n", k, s.Provider.Properties[k])
		}
	}

	// Inputs
	if len(s.Inputs) > 0 {
		fmt.Fprintln(&b, "\n📥 Inputs:")
		keys := sortedKeys(s.Inputs)
		for _, name := range keys {
			input := s.Inputs[name]
			fmt.Fprintf(&b, "  - %s (%s)", name, input.Type)
			if input.Required {
				fmt.Fprint(&b, " [required]")
			}
			if input.Default != nil {
				fmt.Fprintf(&b, " = %v", input.Default)
			}
			if input.Label != "" {
				fmt.Fprintf(&b, "  # %s", input.Label)
			}
			fmt.Fprintln(&b)
			if input.Description != "" {
				fmt.Fprintf(&b, "      → %s\n", input.Description)
			}
			if len(input.Allowed) > 0 {
				fmt.Fprint(&b, "      Allowed:\n")
				for _, v := range input.Allowed {
					fmt.Fprintf(&b, "        - %v (%s)\n", v.Value, v.Label)
				}
			}
		}
	}

	// Secrets
	if len(s.Secrets) > 0 {
		fmt.Fprintln(&b, "\n🔐 Secrets:")
		keys := sortedKeys(s.Secrets)
		for _, name := range keys {
			fmt.Fprintf(&b, "  - %s (type: %s)\n", name, s.Secrets[name])
		}
	}

	// Layers & Steps
	if len(s.Layers) > 0 {
		fmt.Fprintln(&b, "\n🧱 Layers:")
		for i, layer := range s.Layers {
			fmt.Fprintf(&b, "  %d. %s\n", i+1, layer.Name)
			for j, step := range layer.Steps {
				action := step.Action
				if action == "" {
					// fallback to raw action key
					for k := range step.Raw {
						action = k
						break
					}
				}
				fmt.Fprintf(&b, "     %d.%d Step: %s [%s]\n", i+1, j+1, step.Name, action)

				if step.Params != nil {
					fmt.Fprintln(&b, "       Params:")
					paramKeys := sortedKeys(step.Params)
					for _, k := range paramKeys {
						fmt.Fprintf(&b, "         - %s: %v\n", k, step.Params[k])
					}
				}

				if step.Register != "" {
					fmt.Fprintf(&b, "       → registers: %s\n", step.Register)
				}
				if len(step.Tags) > 0 {
					fmt.Fprintf(&b, "       Tags: %v\n", step.Tags)
				}
			}
		}
	}

	// Outputs
	if len(s.Outputs) > 0 {
		fmt.Fprintln(&b, "\n📤 Outputs:")
		keys := sortedKeys(s.Outputs)
		for _, name := range keys {
			out := s.Outputs[name]
			fmt.Fprintf(&b, "  - %s = %s\n", name, out.Value)
			if out.Description != "" {
				fmt.Fprintf(&b, "      → %s\n", out.Description)
			}
		}
	}

	return b.String()
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
