package openapi3spec

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Contact information for the exposed API.
type Contact struct {
	Name  *string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   *string `json:"url,omitempty" yaml:"url,omitempty"`
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate contact object
func (c *Contact) Validate() error {
	if c == nil {
		return nil
	}

	if c.Name != nil && len(strings.TrimSpace(*c.Name)) == 0 {
		return errors.New("info.contact.name if present must not be blank")
	}

	if c.URL != nil {
		_, err := url.Parse(*c.URL)
		if err != nil {
			return fmt.Errorf("info.contact.url if present must be a url: %w", err)
		}
	}

	if c.Email != nil && !rgxEmail.MatchString(*c.Email) {
		return errors.New("info.contact.email if present must be an e-mail address")
	}

	return nil
}
