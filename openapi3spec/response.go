package openapi3spec

import (
	"errors"
	"fmt"
	"strings"
)

// Responses contain possible responses from an Operation
// In order to preserve the data structure we do not allow any extensions
// on Responses
type Responses map[string]*ResponseRef

// Response is a single response from an operation
type Response struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Headers     map[string]*HeaderRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty" yaml:"links,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate response
func (r *Response) Validate() error {
	if len(strings.TrimSpace(r.Description)) == 0 {
		return errors.New("description must not be blank")
	}
	for k, h := range r.Headers {
		if err := h.Validate(); err != nil {
			return fmt.Errorf("headers(%s).%w", k, err)
		}
	}
	for k, c := range r.Content {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("content(%s).%w", k, err)
		}
	}
	for k, l := range r.Links {
		if err := l.Validate(); err != nil {
			return fmt.Errorf("links(%s).%w", k, err)
		}
	}

	if r.Content != nil {
		if _, ok := r.Content["application/json"]; !ok {
			return errors.New("content: must have application/json key if not empty")
		}
	}

	return nil
}

// ResponseRef response reference
type ResponseRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Response
}

// Validate response ref
func (r *ResponseRef) Validate() error {
	// Don't validate references
	if r == nil || len(r.Ref) != 0 {
		return nil
	}

	if err := r.Response.Validate(); err != nil {
		return err
	}

	return nil
}
