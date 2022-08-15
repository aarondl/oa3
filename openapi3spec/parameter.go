package openapi3spec

import (
	"errors"
	"fmt"
	"strings"
)

// Parameter for an operation
type Parameter struct {
	Name            string  `json:"name,omitempty" yaml:"name,omitempty"`
	In              string  `json:"in,omitempty" yaml:"in,omitempty"`
	Description     *string `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool    `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`

	Style         *string    `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       *bool      `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool       `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema        *SchemaRef `json:"schema,omitempty" yaml:"schema,omitempty"`

	Example  any                 `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`

	Content map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate param
func (p *Parameter) Validate(pathTemplates []string) error {
	if len(strings.TrimSpace(p.Name)) == 0 {
		return errors.New("name must not be blank")
	}

	switch p.In {
	case "path":
		found := false
		for _, t := range pathTemplates {
			if t == p.Name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("name %s not found in path templates: [%s]",
				p.Name, strings.Join(pathTemplates, ", "))
		}

		p.Required = true
		if p.AllowEmptyValue {
			return errors.New("allowEmptyValue must not be false for path parameters")
		}
		if p.AllowReserved {
			return errors.New("allowReserved must not be true for path parameters")
		}

		if p.Style == nil {
			p.Style = new(string)
			*p.Style = "simple"
		}
	case "query":
		if p.Style == nil {
			p.Style = new(string)
			*p.Style = "form"
		}
	case "header":
		if p.AllowEmptyValue {
			return errors.New("allowEmptyValue must not be false for header parameters")
		}
		if p.AllowReserved {
			return errors.New("allowReserved must not be true for header parameters")
		}

		if p.Style == nil {
			p.Style = new(string)
			*p.Style = "simple"
		}
	case "cookie":
		if p.AllowEmptyValue {
			return errors.New("allowEmptyValue must not be false for path parameters")
		}
		if p.AllowReserved {
			return errors.New("allowReserved must not be true for cookie parameters")
		}

		if p.Style == nil {
			p.Style = new(string)
			*p.Style = "form"
		}
	default:
		return errors.New("in must be one of: path|query|header|cookie")
	}

	if p.Style != nil && p.Explode == nil {
		p.Explode = new(bool)
		*p.Explode = *p.Style == "form"
	}
	if p.Style != nil {
		switch *p.Style {
		case "matrix", "label":
			if p.In != "path" {
				return fmt.Errorf("style %q can only be used in path", *p.Style)
			}
		case "form":
			if p.In != "query" && p.In != "cookie" {
				return fmt.Errorf("style %q can only be used in query or cookie", *p.Style)
			}
		case "simple":
			if p.In != "path" && p.In != "header" {
				return fmt.Errorf("style %q can only be used in path or header", *p.Style)
			}
			if p.Schema.Schema.Type == "object" {
				return fmt.Errorf("schema can not be of type 'object' when parameter style is simple, got: %s", p.Schema.Schema.Type)
			}
		case "spaceDelimited", "pipeDelimited":
			if p.In != "query" {
				return fmt.Errorf("style %q can only be used in query", *p.Style)
			}
			if p.Schema.Schema.Type != "array" && p.Schema.Schema.Type != "object" {
				return fmt.Errorf("schema must be of type 'object' or 'array' when parameter style is spaceDelimited or pipeDelimited, got: %s", p.Schema.Schema.Type)
			}
		case "deepObject":
			if p.In != "query" {
				return fmt.Errorf("style %q can only be used in query", *p.Style)
			}
			if p.Schema.Schema.Type != "object" {
				return fmt.Errorf("schema must be of type 'object' when parameter style is deepObject, got: %s", p.Schema.Schema.Type)
			}
		default:
			return fmt.Errorf("style must be one of matrix|label|form|simple|spaceDelimited|pipeDelimited|deepObject but found %s", *p.Style)
		}
	}

	if p.Description != nil && len(strings.TrimSpace(*p.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	if p.Schema != nil && len(p.Content) != 0 {
		return errors.New("schema and content are mutually exclusive, define one or the other")
	}

	if p.Schema != nil {
		if err := p.Schema.Validate(); err != nil {
			return fmt.Errorf("schema.%w", err)
		}
	}

	for k, c := range p.Content {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("content(%s).%w", k, err)
		}
	}

	return nil
}

// ParameterRef refers to a parameter
type ParameterRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Parameter
}

// Validate param ref
func (p *ParameterRef) Validate(pathTemplates []string) error {
	// Don't validate references
	if p == nil || len(p.Ref) != 0 {
		return nil
	}

	if err := p.Parameter.Validate(pathTemplates); err != nil {
		return err
	}

	return nil
}

func paramDuplicateKeyCheck(params []*ParameterRef) error {
	if len(params) == 0 {
		return nil
	}

	keys := make(map[string]struct{})
	for _, p := range params {
		_, ok := keys[p.Name+p.In]
		if ok {
			return fmt.Errorf("name %s is duplicated where in is: %s", p.Name, p.In)
		}

		if !p.Required && p.In == "path" {
			return fmt.Errorf(`when in="path" then "required" is itself required and must be set to true`)
		}
	}

	return nil
}
