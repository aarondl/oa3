package openapi3spec

// Parameter for an operation
type Parameter struct {
	Name            string  `json:"name,omitempty" yaml:"name,omitempty"`
	In              string  `json:"in,omitempty" yaml:"in,omitempty"`
	Description     *string `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool    `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`

	Style         *string    `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool       `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool       `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema        *SchemaRef `json:"schema,omitempty" yaml:"schema,omitempty"`

	Example  interface{}         `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`

	Content map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// ParameterRef refers to a parameter
type ParameterRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Parameter
}
