package openapi3spec

// Contact information for the exposed API.
type Contact struct {
	Name  *string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   *string `json:"url,omitempty" yaml:"url,omitempty"`
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
