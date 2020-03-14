package openapi3spec

// Operation on a url
type Operation struct {
	Tags         []string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      *string       `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  *string       `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	OperationID string                  `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  []ParameterRef          `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBodyRef         `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   Responses               `json:"responses,omitempty" yaml:"responses,omitempty"`
	Callbacks   map[string]*CallbackRef `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`

	Deprecated bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security   []SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Servers    []Server              `json:"servers,omitempty" yaml:"servers,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
