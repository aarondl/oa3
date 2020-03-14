package openapi3spec

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to ACL constraints.
//
// Technically Paths can have extensions as per the spec, but we make a choice
// not to conform in order to be able to avoid an object graph that is also
// against the spec: OpenAPI3.Paths.Paths["/url"]
type Paths map[string]*PathRef

// Path describes the operations available on a single path. A Path Item MAY
// be empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available.
type Path struct {
	Summary     string         `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string         `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation     `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation     `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation     `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation     `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation     `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation     `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation     `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation     `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []Server       `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []ParameterRef `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// PathRef refers to a path item
type PathRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Path
}
