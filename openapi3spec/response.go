package openapi3spec

// Responses contain possible responses from an Operation
// In order to preserve the data structure we do not allow any extensions
// on Responses
type Responses map[string]*ResponseRef

// Response is a single response from an operation
type Response struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Headers     map[string]*HeaderRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty" yaml:"links,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// ResponseRef response reference
type ResponseRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Response
}
