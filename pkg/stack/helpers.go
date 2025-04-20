package stack

import (
	"errors"
	"fmt"
	"text/template"
	"text/template/parse"
)

var ErrTemplateSyntax = errors.New("invalid template syntax")

func ExtractVariablePaths(tmplStr string) ([][]string, error) {
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTemplateSyntax, err)
	}

	var variables [][]string

	var walk func(node parse.Node)
	walk = func(node parse.Node) {
		switch n := node.(type) {
		case *parse.ListNode:
			for _, sub := range n.Nodes {
				walk(sub)
			}
		case *parse.ActionNode:
			walk(n.Pipe)
		case *parse.PipeNode:
			for _, cmd := range n.Cmds {
				walk(cmd)
			}
		case *parse.CommandNode:
			for _, arg := range n.Args {
				walk(arg)
			}
		case *parse.VariableNode:
			if len(n.Ident) > 0 {
				// Filter out the dollar sign
				variables = append(variables, n.Ident[1:])
			}
		}
	}

	walk(tmpl.Tree.Root)
	return variables, nil
}
