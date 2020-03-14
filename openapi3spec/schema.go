package openapi3spec

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

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
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

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

// Validate schema
func (s *Schema) Validate(c Components) error {
	if s.Title != nil && len(strings.TrimSpace(*s.Title)) == 0 {
		return errors.New("title if present must not be blank")
	}
	if s.Description != nil && len(strings.TrimSpace(*s.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}
	if len(strings.TrimSpace(s.Type)) == 0 {
		return errors.New("type must not be blank")
	}

	if s.ReadOnly && s.WriteOnly {
		return errors.New("readOnly may not be true at the same time as writeOnly")
	}

	countExclusives := 0
	if len(s.AllOf) != 0 {
		countExclusives++
	}
	if len(s.AnyOf) != 0 {
		countExclusives++
	}
	if len(s.OneOf) != 0 {
		countExclusives++
	}
	if countExclusives > 1 {
		return errors.New("allOf|anyOf|oneOf are mutually exclusive")
	}

	if s.Discriminator != nil {
		if countExclusives == 0 {
			return errors.New("discriminator may only be present with allOf|anyOf|oneOf")
		}

		_, ok := s.Properties[s.Discriminator.PropertyName]
		if !ok {
			return fmt.Errorf("discriminator.propertyName has %s but this property was not found",
				s.Discriminator.PropertyName)
		}
	}

	return nil
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
	Bool bool `json:"bool,omitempty" yaml:"bool,omitempty"`
	*SchemaRef
}

// UnmarshalYAMLObject is called by the unmarshaller to deal with this horrible
// struct.
func (a *AdditionalProperties) UnmarshalYAMLObject(intf interface{}) error {
	if b, ok := intf.(bool); ok {
		a.Bool = b
		return nil
	}

	if m, ok := intf.(map[interface{}]interface{}); ok {
		a.SchemaRef = new(SchemaRef)
		return allocAndSet(reflect.ValueOf(a.SchemaRef).Elem(), m)
	}

	return fmt.Errorf("failed to unmarshal %T into AdditionalProperties", intf)
}

// SchemaRef refers to a schema object
type SchemaRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Schema
}

// Validate schema ref
func (s *SchemaRef) Validate(c Components) error {
	// Don't validate references
	if s == nil || len(s.Ref) != 0 {
		return nil
	}

	if err := s.Schema.Validate(c); err != nil {
		return err
	}

	return nil
}
