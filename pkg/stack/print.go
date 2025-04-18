package stack

import (
	"fmt"
)

// PrintPlan prints a human-readable summary of resources to be created.
func (s *Stack) PrintPlan() {
	fmt.Println("Stack Plan Preview:")
	fmt.Printf("\nStack Name: %s\nDescription: %s\n", s.Name, s.Description)
	fmt.Printf("Provider: %s\n", s.Provider)
	fmt.Printf("\nResources to be created:\n\n")

	for name, res := range s.Resources {
		fmt.Printf("  # %s.%s will be created\n", res.Type, name)
		fmt.Println("  + type:", res.Type)
		printIndentedProperties(res.Properties, "  + ")
		fmt.Println()
	}

	if len(s.Data) > 0 {
		fmt.Printf("Data sources to be queried:\n\n")
		for name, d := range s.Data {
			fmt.Printf("  # data.%s (%s)\n", name, d.Type)
			fmt.Println("  + type:", d.Type)
			printIndentedProperties(d.Properties, "  + ")
			fmt.Println()
		}
	}
}

// Helper to print nested maps/arrays nicely
func printIndentedProperties(val interface{}, indent string) {
	switch v := val.(type) {
	case map[string]interface{}:
		for k, v := range v {
			fmt.Printf("%s%s: ", indent, k)
			switch val := v.(type) {
			case map[string]interface{}, []interface{}:
				fmt.Println()
				printIndentedProperties(val, indent+"  ")
			default:
				fmt.Printf("%v\n", val)
			}
		}
	case []interface{}:
		for _, item := range v {
			switch val := item.(type) {
			case map[string]interface{}, []interface{}:
				fmt.Printf("%s- \n", indent)
				printIndentedProperties(val, indent+"  ")
			default:
				fmt.Printf("%s- %v\n", indent, val)
			}
		}
	default:
		fmt.Printf("%s%v\n", indent, v)
	}
}
