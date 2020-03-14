package openapi3spec

import (
	"errors"
	"net/url"
	"strings"
)

// License information for the exposed API.
type License struct {
	Name string  `json:"name,omitempty" yaml:"name,omitempty"`
	URL  *string `json:"url,omitempty" yaml:"url,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate license
func (l *License) Validate() error {
	if l == nil {
		return nil
	}

	if len(strings.TrimSpace(l.Name)) == 0 {
		return errors.New("info.license.name cannot be blank")
	}

	if l.URL != nil {
		_, err := url.Parse(*l.URL)
		if err != nil {
			return errors.New("info.license.url if present must be a url")
		}
	}

	return nil
}
