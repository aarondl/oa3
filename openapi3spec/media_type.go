package openapi3spec

// MediaType provides schemas and examples for the media type
// identified by its key
type MediaType struct {
	Schema   *SchemaRef           `json:"schema,omitempty" yaml:"schema,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`

	Example  interface{}         `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
