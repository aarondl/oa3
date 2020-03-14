package openapi3spec

// Tag adds metadata to a single tag that is used by the Operation Object. It is
// not mandatory to have a Tag Object per tag defined in the Operation Object
// instances.
type Tag struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  *string       `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}
