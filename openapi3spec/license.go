package openapi3spec

// License information for the exposed API.
type License struct {
	Name string  `json:"name,omitempty" yaml:"name,omitempty"`
	URL  *string `json:"url,omitempty" yaml:"url,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
