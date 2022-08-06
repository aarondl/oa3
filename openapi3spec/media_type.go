package openapi3spec

import "fmt"

// MediaType provides schemas and examples for the media type
// identified by its key
type MediaType struct {
	Schema   SchemaRef            `json:"schema,omitempty" yaml:"schema,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`

	Example  any                 `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate media type
func (m *MediaType) Validate() error {
	if err := m.Schema.Validate(); err != nil {
		return fmt.Errorf("schema.%w", err)
	}

	for k, e := range m.Encoding {
		if err := e.Validate(k, "", ""); err != nil {
			return fmt.Errorf("encoding(%s).%w", k, err)
		}
	}

	return nil
}
