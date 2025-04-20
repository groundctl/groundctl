package stack_test

import (
	"testing"

	"github.com/groundctl/groundctl/pkg/stack"
	"github.com/stretchr/testify/assert"
)

func TestStack_PlanPreview(t *testing.T) {
	t.Run("basic plan output includes key sections", func(t *testing.T) {
		s := &stack.Stack{
			Version:     "1.0.0",
			Name:        "webapp",
			DisplayName: "Web App Stack",
			Description: "Deploys a web server and database",
			Provider: stack.Provider{
				Type: "aws",
				Properties: map[string]any{
					"region":  "us-east-1",
					"profile": "default",
				},
			},
			Inputs: map[string]stack.Input{
				"region": {
					Type:        "string",
					Default:     "us-east-1",
					Required:    true,
					Label:       "Deployment region",
					Description: "AWS region to deploy into",
					Allowed: []stack.AllowedValue{
						{Label: "US East", Value: "us-east-1"},
						{Label: "US West", Value: "us-west-2"},
					},
				},
			},
			Secrets: map[string]stack.Secret{
				"db_password": {
					Type:        "string",
					Label:       "Database password",
					Description: "Password for the PostgreSQL database",
				},
			},
			Layers: []stack.Layer{
				{
					Name: "infrastructure",
					Steps: []stack.Step{
						{
							Name:     "Create VPC",
							Action:   "aws.vpc",
							Params:   map[string]any{"cidr_block": "10.0.0.0/16"},
							Register: "vpc_info",
							Tags:     []string{"network"},
						},
					},
				},
			},
			Outputs: map[string]stack.Output{
				"vpc_id": {
					Value:       "{{ register.vpc_info.id }}",
					Description: "ID of the created VPC",
				},
			},
		}

		output := s.PrintPreview()

		assert.Contains(t, output, "Stack Plan Preview: Web App Stack (version: 1.0.0)")
		assert.Contains(t, output, "🔌 Provider: aws")
		assert.Contains(t, output, "📥 Inputs:")
		assert.Contains(t, output, "region (string) [required] = us-east-1")
		assert.Contains(t, output, "Allowed:")
		assert.Contains(t, output, "🔐 Secrets:")
		assert.Contains(t, output, "db_password")
		assert.Contains(t, output, "🧱 Layers:")
		assert.Contains(t, output, "Create VPC")
		assert.Contains(t, output, "📤 Outputs:")
		assert.Contains(t, output, "vpc_id")
	})

	t.Run("empty stack still renders headers", func(t *testing.T) {
		s := &stack.Stack{
			Name:        "empty",
			DisplayName: "Empty Stack",
			Version:     "0.1.0",
		}

		output := s.PrintPreview()

		assert.Contains(t, output, "Stack Plan Preview: Empty Stack (version: 0.1.0)")
		assert.NotContains(t, output, "Inputs:")
		assert.NotContains(t, output, "Secrets:")
		assert.NotContains(t, output, "Layers:")
		assert.NotContains(t, output, "Outputs:")
	})
}
