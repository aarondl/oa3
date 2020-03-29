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
func (s *Schema) Validate() error {
	if s.Title != nil && len(strings.TrimSpace(*s.Title)) == 0 {
		return errors.New("title if present must not be blank")
	}
	if s.Description != nil && len(strings.TrimSpace(*s.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}
	if len(strings.TrimSpace(s.Type)) == 0 {
		return errors.New("type must not be blank")
	}

	if s.Required != nil {
		if len(s.Required) == 0 {
			return errors.New("required if present must not be empty")
		}

		for i, search := range s.Required {
			found := false
			for j, check := range s.Required {
				if i == j {
					continue
				}
				if search == check {
					found = true
					break
				}
			}

			if found {
				return fmt.Errorf("required has duplicate item: %q", i)
			}
		}
	}

	if s.Properties != nil {
		if len(s.Properties) == 0 {
			return errors.New("properties if present must not be empty")
		}

		for name, prop := range s.Properties {
			// If it's nullable or has a default value, we don't need to
			// validate it's required
			if prop.Schema.Nullable || prop.Schema.Default != nil {
				continue
			}

			if !s.IsRequired(name) {
				return fmt.Errorf("properties(%s): must be required, nullable, or have a default value", name)
			}
		}
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
	if countExclusives != 0 && (len(s.Properties) != 0 || len(s.Required) != 0 || s.AdditionalProperties != nil || s.Items != nil) {
		return errors.New("allOf|anyOf|oneOf cannot be used together with properties|required|additionalProperties|items")
	}

	// There's a lot of validation for allOf/anyOf/oneOf that follows
	if s.Discriminator == nil {
		return nil
	}

	if countExclusives == 0 {
		return errors.New("discriminator may only be present with allOf|anyOf|oneOf")
	}

	if len(s.Discriminator.PropertyName) == 0 {
		return errors.New("discriminator if present must not be empty")
	}

	if len(s.OneOf) != 0 || len(s.AnyOf) != 0 {
		schemas := s.OneOf
		kind := "oneOf"
		if len(s.AnyOf) != 0 {
			schemas = s.AnyOf
			kind = "anyOf"
		}

		for i, s := range schemas {
			if len(s.Ref) == 0 {
				return fmt.Errorf("%s[%d]: must be a $ref when using discriminator", kind, i)
			}

			_, ok := s.Properties[s.Discriminator.PropertyName]
			if !ok {
				return fmt.Errorf("discriminator.propertyName(%s): not found in %s[%d]",
					s.Discriminator.PropertyName,
					kind,
					i)
			}

			if !s.IsRequired(s.Discriminator.PropertyName) {
				return fmt.Errorf("discriminator.propertyName(%s): must be a required property in %s[%d]",
					s.Discriminator.PropertyName,
					kind,
					i)
			}
		}

		dupls := map[string]struct{}{}
		for k, v := range s.Discriminator.Mapping {
			if len(v) == 0 {
				return fmt.Errorf("discriminator.mapping(%s): cannot be empty", k)
			}

			if _, duplicated := dupls[v]; duplicated {
				return fmt.Errorf("discriminator.mapping(%s): duplicates value %s", k, v)
			}
			dupls[v] = struct{}{}

			found := false
			var names []string
			for _, s := range schemas {
				names = append(names, s.Ref)
				if v == s.Ref {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf(`discriminator.mapping(%s): could not find ref %s amongst schemas ("%s")`,
					k, v, strings.Join(names, `", "`),
				)
			}
		}
	} else {
		_, ok := s.Properties[s.Discriminator.PropertyName]
		if !ok {
			return fmt.Errorf("discriminator.propertyName has %s but this property was not found",
				s.Discriminator.PropertyName)
		}

		if !s.IsRequired(s.Discriminator.PropertyName) {
			return fmt.Errorf("discriminator.propertyName(%s): must be a required property",
				s.Discriminator.PropertyName)
		}

		if len(s.Discriminator.Mapping) != 0 {
			return errors.New("discriminator.mapping may not be provided with allOf")
		}
	}

	return nil
}

// IsRequired is a helper to see if a property is required
func (s *Schema) IsRequired(prop string) bool {
	found := false
	for _, check := range s.Required {
		if prop == check {
			found = true
			break
		}
	}

	return found
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
func (s *SchemaRef) Validate() error {
	// Don't validate references
	if s == nil || len(s.Ref) != 0 {
		return nil
	}

	if err := s.Schema.Validate(); err != nil {
		return err
	}

	return nil
}
