package openapi3spec

// ExternalDocs points to external documentation
type ExternalDocs struct {
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string  `json:"url,omitempty" yaml:"url,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
