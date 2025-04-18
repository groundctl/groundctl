package stack

import (
	"fmt"
	"regexp"
	"strings"
)

func (s *Stack) ValidateDependsOn() error {
	// Track cyclic references
	visited := map[string]bool{}
	stackPath := map[string]bool{}

	var visit func(name string) error
	visit = func(name string) error {
		if stackPath[name] {
			return fmt.Errorf("cyclic dependency detected involving resource %q", name)
		}
		if visited[name] {
			return nil
		}

		visited[name] = true
		stackPath[name] = true

		res, ok := s.Resources[name]
		if !ok {
			return nil
		}

		for _, dep := range res.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}

		stackPath[name] = false
		return nil
	}

	for resName, res := range s.Resources {
		for _, dep := range res.DependsOn {
			_, exists := s.Resources[dep]
			if !exists {
				return fmt.Errorf("resource %q depends on non-existent resource %q", resName, dep)
			}
		}

		if err := visit(resName); err != nil {
			return err
		}
	}
	return nil
}

// Regex pattern to capture references like ${ec2_instance.public_ip} or ${var.environment}
var referencePattern = regexp.MustCompile(`\$\{([^}]*)\}`)

func (s *Stack) ValidateReferences() error {
	// Make a list of all possible reference keys
	allKeys := make(map[string]struct{})
	for k := range s.Resources {
		allKeys["resource."+k] = struct{}{}
	}
	for k := range s.Data {
		allKeys["data."+k] = struct{}{}
	}
	for k := range s.Variables {
		allKeys["var."+k] = struct{}{}
	}

	// Build reference graph
	graph := make(map[string][]string)

	// Traverse and build graph from all resources, data, outputs
	for name, r := range s.Resources {
		node := "resource." + name
		if err := walkWithGraph(r.Properties, node, allKeys, graph); err != nil {
			return err
		}
	}
	for name, d := range s.Data {
		node := "data." + name
		if err := walkWithGraph(d.Properties, node, allKeys, graph); err != nil {
			return err
		}
	}
	for name, o := range s.Outputs {
		node := "output." + name
		if err := walkWithGraph(o.Value, node, allKeys, graph); err != nil {
			return err
		}
	}

	// Detect cycles in reference graph
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	var visit func(string) error
	visit = func(n string) error {
		if stack[n] {
			return fmt.Errorf("cyclic reference detected involving %q", n)
		}
		if visited[n] {
			return nil
		}
		visited[n] = true
		stack[n] = true

		for _, dep := range graph[n] {
			if err := visit(dep); err != nil {
				return err
			}
		}

		stack[n] = false
		return nil
	}

	for n := range graph {
		if err := visit(n); err != nil {
			return err
		}
	}

	return nil
}

// walkWithGraph validates references and builds a dependency graph
func walkWithGraph(val any, current string, allKeys map[string]struct{}, graph map[string][]string) error {
	switch v := val.(type) {
	case map[string]any:
		for _, v := range v {
			if err := walkWithGraph(v, current, allKeys, graph); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range v {
			if err := walkWithGraph(item, current, allKeys, graph); err != nil {
				return err
			}
		}
	default:
		valStr := fmt.Sprint(v)
		matches := referencePattern.FindAllStringSubmatch(valStr, -1)
		for _, match := range matches {
			if len(match) < 2 || match[1] == "" {
				return fmt.Errorf("invalid reference %s", match[0])
			}
			ref := match[1]
			parts := strings.Split(ref, ".")
			if len(parts) < 2 {
				return fmt.Errorf("invalid reference %s", match[0])
			}
			prefix := parts[0]
			key := parts[1]

			qualified := prefix + "." + key
			if _, ok := allKeys[qualified]; !ok {
				return fmt.Errorf("invalid reference %s", match[0])
			}

			// Build graph edge
			graph[current] = append(graph[current], qualified)
		}
	}
	return nil
}
