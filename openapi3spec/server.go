package openapi3spec

// Server represents a server of this API
type Server struct {
	URL         string                    `json:"url,omitempty" yaml:"url,omitempty"`
	Description *string                   `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// ServerVariable for server URL template substitution
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default,omitempty" yaml:"default,omitempty"`
	Description *string  `json:"description,omitempty" yaml:"description,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
