package rabbitmq

import "github.com/groundctl/groundctl/pkg/stack"

type DeployRequest struct {
	Stack   stack.Stack    `json:"stack"`
	Secrets map[string]any `json:"secrets"`
	Inputs  map[string]any `json:"inputs"`
}
