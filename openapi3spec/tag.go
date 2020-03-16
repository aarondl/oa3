package openapi3spec

import (
	"errors"
	"fmt"
	"strings"
)

// Tag adds metadata to a single tag that is used by the Operation Object. It is
// not mandatory to have a Tag Object per tag defined in the Operation Object
// instances.
type Tag struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  *string       `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Validate a tag
func (t *Tag) Validate() error {
	if len(strings.TrimSpace(t.Name)) == 0 {
		return errors.New("name cannot be blank")
	}

	if t.Description != nil && len(strings.TrimSpace(*t.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	if err := t.ExternalDocs.Validate(); err != nil {
		return fmt.Errorf("externalDocs.%w", err)
	}

	return nil
}
