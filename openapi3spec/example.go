package openapi3spec

// Example object
type Example struct {
	Summary     *string `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// Value and ExternalValue are mutually exclusive
	Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// ExampleRef refers to an example object
type ExampleRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Example
}
