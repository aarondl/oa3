package openapi3spec

import (
	"errors"
	"fmt"
	"strings"
)

// Encoding definition applied to a single schema object
type Encoding struct {
	ContentType *string            `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers     map[string]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`

	Style         *string `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool    `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool    `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate encoding
func (e *Encoding) Validate(mediaType, kind, format string) error {
	if e.ContentType == nil {
		switch {
		case kind == "string" && format == "binary":
			e.ContentType = new(string)
			*e.ContentType = "application/octet-stream"
		case kind == "integer" || kind == "boolean" || kind == "float":
			e.ContentType = new(string)
			*e.ContentType = "text/plain"
		case kind == "object":
			e.ContentType = new(string)
			*e.ContentType = "application/json"
		}
	} else if len(strings.TrimSpace(*e.ContentType)) == 0 {
		return fmt.Errorf("contentType if present must not be blank")
	}

	for k, h := range e.Headers {
		if k == "Content-Disposition" || k == "Content-Type" {
			delete(e.Headers, k)
			continue
		}

		if err := h.Validate(); err != nil {
			return fmt.Errorf("headers(%s).%w", k, err)
		}
	}

	if e.Style != nil && *e.Style == "form" {
		e.Explode = true
	}
	if e.Style != nil {
		switch *e.Style {
		case "matrix", "label", "form", "simple", "spaceDelimited", "pipeDelimited", "deepObject":
			return errors.New("style must be one of matrix|label|form|simple|spaceDelimited|pipeDelimited|deepObject")
		}
	}

	if mediaType != "application/x-www-form-urlencoded" {
		e.AllowReserved = false
	}

	return nil
}
