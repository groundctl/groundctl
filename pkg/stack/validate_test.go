package stack_test

import (
	"testing"

	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/stretchr/testify/assert"
)

func TestStackValidate(t *testing.T) {
	t.Run("missing metadata", func(t *testing.T) {
		s := &stack.Stack{
			Version: "",
			Name:    "",
			Provider: stack.Provider{
				Type: "",
			},
		}
		err := s.Validate()
		assert.ErrorContains(t, err, "missing required metadata fields")
	})

	t.Run("invalid input fields", func(t *testing.T) {
		s := &stack.Stack{
			Version: "1",
			Name:    "test",
			Provider: stack.Provider{
				Type: "mock",
			},
			Inputs: map[string]stack.Input{
				"region": {Type: ""},
			},
		}
		err := s.Validate()
		assert.ErrorContains(t, err, "input 'region' must have a type")
	})

	t.Run("invalid input default", func(t *testing.T) {
		s := &stack.Stack{
			Version: "1",
			Name:    "test",
			Provider: stack.Provider{
				Type: "aws",
			},
			Inputs: map[string]stack.Input{
				"env": {
					Type:     "string",
					Required: false,
					Default:  nil,
				},
			},
		}
		err := s.Validate()
		assert.ErrorContains(t, err, "input 'env' is not required but has no default value")
	})

	t.Run("input allow list check", func(t *testing.T) {
		s := &stack.Stack{
			Version: "1",
			Name:    "test",
			Provider: stack.Provider{
				Type: "aws",
			},
			Inputs: map[string]stack.Input{
				"region": {
					Type:     "string",
					Required: false,
					Default:  "us-west-2",
					Allowed: []stack.AllowedValue{
						{Value: "us-east-1"},
						{Value: "eu-west-1"},
					},
				},
			},
		}
		err := s.Validate()
		assert.ErrorContains(t, err, "default value of input 'region' is not in the allowed list")
	})

	t.Run("input template with valid reference", func(t *testing.T) {
		s := &stack.Stack{
			Version: "1",
			Name:    "test",
			Provider: stack.Provider{
				Type:       "aws",
				Properties: map[string]any{"region": "{{ $.input.region }}"},
			},
			Inputs: map[string]stack.Input{
				"region": {Type: "string", Required: true},
			},
			Secrets:             map[string]stack.Secret{},
			RegisteredVariables: map[string]map[string]any{},
		}
		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("input template with undefined reference", func(t *testing.T) {
		s := &stack.Stack{
			Version: "1",
			Name:    "Test",
			Provider: stack.Provider{
				Type:       "aws",
				Properties: map[string]any{"region": "{{ $.input.env }}"},
			},
			Inputs: map[string]stack.Input{
				"region": {Type: "string", Required: true},
			},
			RegisteredVariables: map[string]map[string]any{},
		}
		err := s.Validate()
		assert.ErrorContains(t, err, "undefined input 'env'")
	})
}
