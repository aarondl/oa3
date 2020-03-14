package openapi3spec

// Components specify referenceable reusable components
type Components struct {
	Schemas         map[string]*SchemaRef         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*ResponseRef       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameter       map[string]*ParameterRef      `json:"parameter,omitempty" yaml:"parameter,omitempty"`
	Examples        map[string]*ExampleRef        `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBodyRef    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]*HeaderRef         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeRef `json:"securityScheme,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]*LinkRef           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]*CallbackRef       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}
