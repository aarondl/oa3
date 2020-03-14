package openapi3spec

// RequestBody for an operation
type RequestBody struct {
	Description *string               `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`
}

// RequestBodyRef refers to a request body
type RequestBodyRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*RequestBody
}
