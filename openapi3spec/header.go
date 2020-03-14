package openapi3spec

// Header object
type Header struct {
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

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// HeaderRef refers to a parameter
type HeaderRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Header
}
