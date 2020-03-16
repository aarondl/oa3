package openapi3spec

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ExternalDocs points to external documentation
type ExternalDocs struct {
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string  `json:"url,omitempty" yaml:"url,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate external docs
func (e *ExternalDocs) Validate() error {
	if e.Description != nil && len(strings.TrimSpace(*e.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	_, err := url.Parse(e.URL)
	if err != nil {
		return fmt.Errorf("url must be a valid url: %w", err)
	}

	return nil
}
