package openapi3spec

// Schema Object allows the definition of input and output data types. These
// types can be objects, but also primitives and arrays. This object is an
// extended subset of the JSON Schema Specification Wright Draft 00. For more
// information about the properties, see JSON Schema Core and JSON Schema
// Validation. Unless stated otherwise, the property definitions follow the JSON
// Schema.
type Schema struct {
	Title       *string     `json:"title,omitempty" yaml:"title,omitempty"`
	Description *string     `json:"description,omitempty" yaml:"description,omitempty"`
	Default     interface{} `json:"default,omitempty" yaml:"default,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	Nullable   bool `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	ReadOnly   bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly  bool `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	Example      interface{}   `json:"example,omitempty" yaml:"example,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	MultipleOf       *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

	MaxLength *int `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength *int `json:"minLength,omitempty" yaml:"minLength,omitempty"`

	Format  *string `json:"format,omitempty" yaml:"format,omitempty"`
	Pattern *string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	Items       *SchemaRef `json:"items,omitempty" yaml:"items,omitempty"`
	MaxItems    *int       `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems    *int       `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems *bool      `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`

	Required      []string      `json:"required,omitempty" yaml:"required,omitempty"`
	Enum          []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`
	MaxProperties *int          `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties *int          `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`

	Properties           map[string]*SchemaRef `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`

	AllOf         []*SchemaRef   `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf         []*SchemaRef   `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	OneOf         []*SchemaRef   `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Not           *SchemaRef     `json:"not,omitempty" yaml:"not,omitempty"`
	Discriminator *Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Discriminator helps with decoding. When request bodies or response payloads
// may be one of a number of different schemas, a discriminator object can be
// used to aid in serialization, deserialization, and validation. The
// discriminator is a specific object in a schema which is used to inform the
// consumer of the specification of an alternative schema based on the value
// associated with it.
type Discriminator struct {
	PropertyName string            `json:"propertyName,omitempty" yaml:"propertyName,omitempty"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

// AdditionalProperties is ridiculous, a bool or a schema
type AdditionalProperties struct {
	Bool   bool       `json:"bool,omitempty" yaml:"bool,omitempty"`
	Schema *SchemaRef `json:"schema,omitempty" yaml:"schema,omitempty"`
}

// SchemaRef refers to a schema object
type SchemaRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Schema
}
