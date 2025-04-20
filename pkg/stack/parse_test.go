package stack_test

import (
	"testing"

	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStackParse(t *testing.T) {
	t.Run("parse valid stack template", func(t *testing.T) {
		yamlData := `
version: "1.0"
name: "sample"
display_name: "Sample"
description: "Test stack"
provider:
  type: git@github.com:example/provider
layers:
  - name: setup
    steps:
      - name: do something
        test.action:
          key: value
`
		stk, err := stack.Parse([]byte(yamlData))
		require.NoError(t, err)
		assert.Equal(t, "sample", stk.Name)
		assert.Equal(t, "test.action", stk.Layers[0].Steps[0].Action)
		assert.Equal(t, map[string]any{"key": "value"}, stk.Layers[0].Steps[0].Params)

	})

	t.Run("parse invalid stack template", func(t *testing.T) {
		yamlData := `
version: "1.0"
name: "bad"
provider:
  type: example
layers:
  - name: badlayer
    steps:
      - name: invalid
        test.action: string_instead_of_map
`
		_, err := stack.Parse([]byte(yamlData))
		assert.Error(t, err)
	})
}
