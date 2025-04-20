package stack

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Resolves all references to given inputs and secrets to prepare for a deployment
func (s *Stack) Resolve(inputs map[string]any, secrets map[string]string) error {
	ctx := map[string]any{
		"input":  inputs,
		"secret": secrets,
	}

	var tmplEval func(val any) (any, error)
	tmplEval = func(val any) (any, error) {
		switch v := val.(type) {
		case string:
			if !strings.Contains(v, "{{") {
				return v, nil
			}
			tmpl, err := template.New("").Parse(v)
			if err != nil {
				return nil, err
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, ctx); err != nil {
				return nil, err
			}
			return buf.String(), nil
		case map[string]any:
			out := make(map[string]any)
			for k, val := range v {
				res, err := tmplEval(val)
				if err != nil {
					return nil, err
				}
				out[k] = res
			}
			return out, nil
		case []any:
			for i := range v {
				res, err := tmplEval(v[i])
				if err != nil {
					return nil, err
				}
				v[i] = res
			}
			return v, nil
		default:
			return v, nil
		}
	}

	providerProps, err := tmplEval(s.Provider.Properties)
	if err != nil {
		return fmt.Errorf("failed to resolve provider properties: %w", err)
	}
	s.Provider.Properties = providerProps.(map[string]any)

	for i := range s.Layers {
		for j := range s.Layers[i].Steps {
			params, err := tmplEval(s.Layers[i].Steps[j].Params)
			if err != nil {
				return fmt.Errorf("failed to resolve step '%s' in layer '%s': %w", s.Layers[i].Steps[j].Name, s.Layers[i].Name, err)
			}
			s.Layers[i].Steps[j].Params = params.(map[string]any)
		}
	}

	return nil
}
