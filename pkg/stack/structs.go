package stack

type Stack struct {
	Version     string            `yaml:"version"`
	Name        string            `yaml:"name"`
	DisplayName string            `yaml:"display_name"`
	Description string            `yaml:"description"`
	Provider    Provider          `yaml:"provider"`
	Secrets     map[string]string `yaml:"secrets,omitempty"`
	Inputs      map[string]Input  `yaml:"inputs,omitempty"`
	Layers      []Layer           `yaml:"layers"`
	Outputs     map[string]Output `yaml:"outputs"`
	// Contains all registered variables from steps that contain a "register" attribute
	RegisteredVariables map[string]map[string]any
}

type Provider struct {
	Type       string         `yaml:"type"`
	Properties map[string]any `yaml:"properties,omitempty"`
}

type Secret struct {
	Type        string         `yaml:"type"`
	Allowed     []AllowedValue `yaml:"allowed,omitempty"`
	Description string         `yaml:"description,omitempty"`
	Label       string         `yaml:"label,omitempty"`
}

type Input struct {
	Type        string         `yaml:"type"`
	Default     any            `yaml:"default"`
	Allowed     []AllowedValue `yaml:"allowed,omitempty"`
	Required    bool           `yaml:"required,omitempty"`
	Description string         `yaml:"description,omitempty"`
	Label       string         `yaml:"label,omitempty"`
}

type AllowedValue struct {
	Label string `yaml:"label"`
	Value any    `yaml:"value"`
}

type Layer struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name     string         `yaml:"name"`
	Action   string         `yaml:"-"`
	Params   map[string]any `yaml:"-"`
	Register string         `yaml:"register,omitempty"`
	Tags     []string       `yaml:"tags,omitempty"`
	Raw      map[string]any `yaml:",inline"`
}

type Output struct {
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}
