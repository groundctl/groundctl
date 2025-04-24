package stack

type Stack struct {
	Version     string            `yaml:"version" json:"version"`
	Name        string            `yaml:"name" json:"name"`
	DisplayName string            `yaml:"display_name" json:"display_name"`
	Description string            `yaml:"description" json:"description"`
	Provider    Provider          `yaml:"provider" json:"provider"`
	Secrets     map[string]Secret `yaml:"secrets,omitempty" json:"secrets,omitempty"`
	Inputs      map[string]Input  `yaml:"inputs,omitempty" json:"inputs,omitempty"`
	Layers      []Layer           `yaml:"layers" json:"layers"`
	Outputs     map[string]Output `yaml:"outputs" json:"outputs"`
	// Contains all registered variables from steps that contain a "register" attribute
	RegisteredVariables map[string]map[string]any `json:"registered_variables"`
}

type Provider struct {
	Type       string         `yaml:"type" json:"type"`
	Properties map[string]any `yaml:"properties,omitempty" json:"properties,omitempty"`
}

type Secret struct {
	Type        string         `yaml:"type" json:"type"`
	Allowed     []AllowedValue `yaml:"allowed,omitempty" json:"allowed,omitempty"`
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	Label       string         `yaml:"label,omitempty" json:"label,omitempty"`
}

type Input struct {
	Type        string         `yaml:"type" json:"type"`
	Default     any            `yaml:"default" json:"default"`
	Allowed     []AllowedValue `yaml:"allowed,omitempty" json:"allowed,omitempty"`
	Required    bool           `yaml:"required,omitempty" json:"required,omitempty"`
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	Label       string         `yaml:"label,omitempty" json:"label,omitempty"`
}

type AllowedValue struct {
	Label string `yaml:"label" json:"label"`
	Value any    `yaml:"value" json:"value"`
}

type Layer struct {
	Name  string `yaml:"name" json:"name"`
	Steps []Step `yaml:"steps" json:"steps"`
}

type Step struct {
	Name     string         `yaml:"name" json:"name"`
	Action   string         `yaml:"-" json:"-"`
	Params   map[string]any `yaml:"-" json:"-"`
	Register string         `yaml:"register,omitempty" json:"register,omitempty"`
	Tags     []string       `yaml:"tags,omitempty" json:"tags,omitempty"`
	Raw      map[string]any `yaml:",inline" json:",inline"`
}

type Output struct {
	Value       string `yaml:"value" json:"value"`
	Description string `yaml:"description" json:"description"`
}
