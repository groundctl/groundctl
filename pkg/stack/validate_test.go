package stack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDependsOn(t *testing.T) {
	t.Run("invalid resource depends_on array", func(t *testing.T) {
		stack := Stack{
			Resources: map[string]StackResource{},
		}
		stack.Resources["ec2_instance"] = StackResource{
			DependsOn: []string{"non_existent_resource"},
		}
		err := stack.ValidateDependsOn()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "depends on non-existent resource")
	})

	t.Run("cyclic depends_on reference", func(t *testing.T) {
		stack := Stack{
			Resources: map[string]StackResource{},
		}
		stack.Resources["r1"] = StackResource{
			Type:      "aws/instance",
			DependsOn: []string{"r2"},
		}
		stack.Resources["r2"] = StackResource{
			Type:      "aws/instance",
			DependsOn: []string{"r3"},
		}
		stack.Resources["r3"] = StackResource{
			Type:      "aws/instance",
			DependsOn: []string{"r1"}, // Cycle here
		}
		err := stack.ValidateDependsOn()
		assert.Error(t, err)
		fmt.Printf("%v\n", err)
		assert.Contains(t, err.Error(), "cyclic dependency")
	})
}

func TestValidateReferences(t *testing.T) {
	t.Run("validate valid references", func(t *testing.T) {
		stack := Stack{
			Variables: map[string]StackVariable{
				"ami_id": {
					Type:    VarTypeString,
					Default: "ubuntu2204",
				},
			},
			Resources: map[string]StackResource{
				"my_ec2": {
					Type: "aws/ec2",
					Properties: map[string]any{
						"ami": "${var.ami_id}",
					},
				},
			},
		}
		err := stack.ValidateReferences()
		assert.NoError(t, err)
	})

	t.Run("invalid data reference", func(t *testing.T) {
		fmt.Println("invalid data reference test")
		stack := Stack{
			Data: map[string]StackData{
				"my_network": {
					Type: "aws/network",
				},
			},
			Resources: map[string]StackResource{
				"my_ec2": {
					Type: "aws/ec2",
					Properties: map[string]any{
						"interfaces": []struct {
							Network string
						}{
							{
								Network: "${data.non_existent_network.id}",
							},
						},
					},
				},
			},
		}
		err := stack.ValidateReferences()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reference ${data.non_existent_network.id}")
	})

	t.Run("invalid variable reference", func(t *testing.T) {
		stack := Stack{
			Variables: map[string]StackVariable{
				"ami_id": {
					Type:    VarTypeString,
					Default: "ubuntu2204",
				},
			},
			Resources: map[string]StackResource{
				"my_ec2": {
					Type: "aws/ec2",
					Properties: map[string]any{
						"ami": "${var.non_existent_variable}",
					},
				},
			},
		}
		err := stack.ValidateReferences()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reference ${var.non_existent_variable}")
	})

	t.Run("empty reference", func(t *testing.T) {
		stack := Stack{
			Resources: map[string]StackResource{
				"my_ec2": {
					Type: "aws/ec2",
					Properties: map[string]any{
						"ami": "${}",
					},
				},
			},
		}
		err := stack.ValidateReferences()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reference")
	})

	t.Run("valid complex reference", func(t *testing.T) {
		stack := Stack{
			Data: map[string]StackData{
				"my_network": {
					Type: "aws/network",
				},
			},
			Resources: map[string]StackResource{
				"my_ec2": {
					Type: "aws/ec2",
					Properties: map[string]any{
						"interfaces": []struct {
							Network string
						}{
							{
								Network: "${data.my_network.id}",
							},
						},
					},
				},
			},
		}
		err := stack.ValidateReferences()
		assert.NoError(t, err)
	})

	t.Run("cyclic references", func(t *testing.T) {
		stack := Stack{
			Resources: map[string]StackResource{},
		}
		stack.Resources["one"] = StackResource{
			Type: "aws/instance",
			Properties: map[string]any{
				"value": "${resource.two.id}",
			},
		}
		stack.Resources["two"] = StackResource{
			Type: "aws/instance",
			Properties: map[string]any{
				"value": "${resource.one.id}",
			},
		}
		err := stack.ValidateReferences()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cyclic reference")
	})
}
