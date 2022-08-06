package openapi3spec

import (
	"errors"
	"fmt"
	"strings"
)

// Header object
type Header struct {
	Description     *string `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool    `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`

	Style         *string    `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool       `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool       `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema        *SchemaRef `json:"schema,omitempty" yaml:"schema,omitempty"`

	Example  any                 `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`

	Content map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate header
func (h *Header) Validate() error {
	if h.AllowEmptyValue {
		return errors.New("allowEmptyValue must not be false for header parameters")
	}
	if h.AllowReserved {
		return errors.New("allowReserved must not be true for header parameters")
	}

	if h.Style == nil {
		h.Style = new(string)
		*h.Style = "simple"
	}

	if h.Style != nil && *h.Style == "form" {
		h.Explode = true
	}
	if h.Style != nil {
		switch *h.Style {
		case "matrix", "label", "form", "simple", "spaceDelimited", "pipeDelimited", "deepObject":
		default:
			return errors.New("style must be one of matrix|label|form|simple|spaceDelimited|pipeDelimited|deepObject")
		}
	}

	if h.Description != nil && len(strings.TrimSpace(*h.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	if err := h.Schema.Validate(); err != nil {
		return fmt.Errorf("schema.%w", err)
	}

	for k, c := range h.Content {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("content(%s).%w", k, err)
		}
	}

	return nil
}

// HeaderRef refers to a parameter
type HeaderRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Header
}

// Validate header ref
func (h *HeaderRef) Validate() error {
	// Don't validate references
	if h == nil || len(h.Ref) != 0 {
		return nil
	}

	if err := h.Header.Validate(); err != nil {
		return err
	}

	return nil
}
