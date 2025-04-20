package stack_test

import (
	"testing"

	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/stretchr/testify/assert"
)

func TestExtractVariablePaths(t *testing.T) {
	t.Run("valid templates", func(t *testing.T) {
		cases := []struct {
			name     string
			input    string
			expected [][]string
		}{
			{
				name:     "single variable",
				input:    `{{ $.input.region }}`,
				expected: [][]string{{"input", "region"}},
			},
			{
				name:     "multiple variables",
				input:    `{{ $.input.env }}-{{ $.input.region }}`,
				expected: [][]string{{"input", "env"}, {"input", "region"}},
			},
			{
				name:     "custom variable",
				input:    `{{ $.my_var.prop }}`,
				expected: [][]string{{"my_var", "prop"}},
			},
		}
		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				paths, err := stack.ExtractVariablePaths(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, paths)
			})
		}
	})

	t.Run("invalid template", func(t *testing.T) {
		_, err := stack.ExtractVariablePaths("{{ .input.region ")
		assert.ErrorIs(t, err, stack.ErrTemplateSyntax)
	})
}
