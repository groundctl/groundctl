package stack

type Stack struct {
	Name        string                   `yaml:"name"`
	DisplayName string                   `yaml:"display_name,omitempty"`
	Description string                   `yaml:"description"`
	Provider    string                   `yaml:"provider"`
	Variables   map[string]StackVariable `yaml:"variables"`
	Data        map[string]StackData     `yaml:"data"`
	Resources   map[string]StackResource `yaml:"resources,omitempty"`
	Outputs     map[string]StackOutput   `yaml:"outputs,omitempty"`
	Imports     []string                 `yaml:"imports,omitempty"`
}

type StackVarType string

const (
	VarTypeString StackVarType = "string"
	VarTypeInt    StackVarType = "int"
	VarTypeFloat  StackVarType = "float"
	VarTypeIP     StackVarType = "ip"
)

type StackVariable struct {
	Type    StackVarType  `yaml:"type"`
	Default string        `yaml:"default"`
	Allowed []interface{} `yaml:"allowed,omitempty"`
}

type StackData struct {
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties"`
}

type StackResource struct {
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties"`
	DependsOn  []string               `yaml:"depends_on"`
}

type StackOutput struct {
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}
