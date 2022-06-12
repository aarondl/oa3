package openapi3spec

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
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
	switch s.Type {
	case "object", "array", "boolean", "number", "integer", "string":
	default:
		return fmt.Errorf("type must be one of object|array|boolean|number|integer|string but got %q", s.Type)
	}
	if len(s.Enum) != 0 && s.Type != "string" {
		return errors.New("enum cannot contain non-strings")
	}
	if s.Type == "array" && s.Items == nil {
		return errors.New("items must be present if type is array")
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
	}

	if s.MultipleOf != nil {
		switch s.Type {
		case "integer", "number":
		default:
			return errors.New("multipleOf: cannot be used unless type is one of: 'integer', 'number'")
		}
	}
	if s.Maximum != nil {
		switch s.Type {
		case "integer", "number":
		default:
			return errors.New("maximum: cannot be used unless type is one of: 'integer', 'number'")
		}
	}
	if s.Minimum != nil {
		switch s.Type {
		case "integer", "number":
		default:
			return errors.New("minimum: cannot be used unless type is one of: 'integer', 'number'")
		}
	}
	if s.Minimum != nil && s.Maximum != nil {
		if *s.Minimum > *s.Maximum {
			return fmt.Errorf("maximum(%f): cannot be less than minimum (was %f)", *s.Maximum, *s.Minimum)
		}
	}
	if s.MaxLength != nil {
		if *s.MaxLength <= 0 {
			return fmt.Errorf("maxLength: must be greater than 0 (was %d)", *s.MaxLength)
		}

		switch s.Type {
		case "string":
		default:
			return errors.New("maxLength: cannot be used unless type is one of: 'string'")
		}
	}
	if s.MinLength != nil {
		if *s.MinLength < 0 {
			return fmt.Errorf("minLength: cannot be a negative number (was %d)", *s.MinLength)
		}

		switch s.Type {
		case "string":
		default:
			return errors.New("minLength: cannot be used unless type is one of: 'string'")
		}
	}
	if s.MinLength != nil && s.MaxLength != nil {
		if *s.MinLength > *s.MaxLength {
			return fmt.Errorf("maxLength(%d): cannot be less than minLength(%d)", *s.MaxLength, *s.MinLength)
		}
	}
	if s.MaxItems != nil {
		if *s.MaxItems <= 0 {
			return fmt.Errorf("maxItems: must be greater than 0 (was %d)", *s.MaxItems)
		}

		switch s.Type {
		case "array":
		default:
			return errors.New("maxItems: cannot be used unless type is one of: 'array'")
		}
	}
	if s.MinItems != nil {
		if *s.MinItems < 0 {
			return fmt.Errorf("minItems: cannot be a negative number (was %d)", *s.MinItems)
		}

		switch s.Type {
		case "array":
		default:
			return errors.New("minItems: cannot be used unless type is one of: 'array'")
		}
	}
	if s.MinItems != nil && s.MaxItems != nil {
		if *s.MinItems > *s.MaxItems {
			return fmt.Errorf("maxItems(%d): cannot be less than minItems(%d)", *s.MaxItems, *s.MinItems)
		}
	}
	if s.UniqueItems != nil {
		switch s.Type {
		case "array":
		default:
			return errors.New("uniqueItems: cannot be used unless type is one of: 'array'")
		}
	}
	if s.MaxProperties != nil {
		if *s.MaxProperties <= 0 {
			return fmt.Errorf("maxProperties: must be greater than 0 (was %d)", *s.MaxProperties)
		}
		if s.AdditionalProperties == nil {
			return fmt.Errorf("maxProperties: cannot use unless additionalProperties is specified")
		}

		switch s.Type {
		case "object":
		default:
			return errors.New("maxProperties: cannot be used unless type is one of: 'object'")
		}
	}
	if s.MinProperties != nil {
		if *s.MinProperties < 0 {
			return fmt.Errorf("minProperties: cannot be a negative number (was %d)", *s.MinProperties)
		}
		if s.AdditionalProperties == nil {
			return fmt.Errorf("minProperties: cannot use unless additionalProperties is specified")
		}

		switch s.Type {
		case "object":
		default:
			return errors.New("minProperties: cannot be used unless type is one of: 'object'")
		}
	}
	if s.MinProperties != nil && s.MaxProperties != nil {
		if *s.MinProperties > *s.MaxProperties {
			return fmt.Errorf("maxProperties(%d): cannot be less than minProperties(%d)", *s.MaxProperties, *s.MinProperties)
		}
	}

	if s.Pattern != nil {
		_, err := regexp.Compile(*s.Pattern)
		if err != nil {
			return fmt.Errorf("pattern(%s): failed to compile regular expression: %w", *s.Pattern, err)
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
		return errors.New("discriminator may only be present with anyOf|oneOf")
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
