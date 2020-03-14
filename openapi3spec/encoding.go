package openapi3spec

// Encoding definition applied to a single schema object
type Encoding struct {
	ContentType *string            `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers     map[string]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`

	Style         *string `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool    `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool    `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
