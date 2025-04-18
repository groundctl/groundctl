package stack

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit test to parse example_stack.yml
func TestParseExampleStack(t *testing.T) {
	data, err := os.ReadFile("example_stack.yml")
	assert.NoError(t, err)

	stack, err := Parse(data)
	assert.NoError(t, err)

	t.Run("parse template", func(t *testing.T) {
		// Assert general fields
		assert.Equal(t, "example-stack", stack.Name)
		assert.Equal(t, "Example Stack", stack.DisplayName)
		assert.Equal(t, "Deploys a simple EC2 instance with security group on AWS", stack.Description)
		assert.Equal(t, "my-aws-satellite", stack.Provider)

		// Check variable parsing
		instanceType, ok := stack.Variables["instance_type"]
		assert.True(t, ok)
		assert.Equal(t, VarTypeString, instanceType.Type)
		assert.Equal(t, "t2.micro", instanceType.Default)

		env, ok := stack.Variables["environment"]
		assert.True(t, ok)
		assert.Equal(t, VarTypeString, env.Type)
		assert.Equal(t, "dev", env.Default)
		assert.ElementsMatch(t, []interface{}{"dev", "staging", "prod"}, env.Allowed)

		// Check resource parsing
		ec2, ok := stack.Resources["ec2_instance"]
		assert.True(t, ok)
		assert.Equal(t, "aws/instance", ec2.Type)
		assert.NotEmpty(t, ec2.Properties["instance_type"])

		// Check the validation for depends_on
		err = stack.ValidateReferences()
		assert.NoError(t, err) // should pass, no missing dependencies
	})
}
