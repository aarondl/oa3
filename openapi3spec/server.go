package openapi3spec

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Server represents a server of this API
type Server struct {
	URL         string                    `json:"url,omitempty" yaml:"url,omitempty"`
	Description *string                   `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate server
func (s *Server) Validate() error {
	if s == nil {
		return nil
	}

	if _, err := url.Parse(s.URL); err != nil {
		return errors.New("url must be a url")
	}

	if s.Description != nil && len(strings.TrimSpace(*s.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	for k, v := range s.Variables {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("serverVariable[%s].%w", k, err)
		}
	}

	return nil
}

// ServerVariable for server URL template substitution
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default,omitempty" yaml:"default,omitempty"`
	Description *string  `json:"description,omitempty" yaml:"description,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate server variable
func (s *ServerVariable) Validate() error {
	if s == nil {
		return nil
	}

	for i, e := range s.Enum {
		if len(strings.TrimSpace(e)) == 0 {
			return fmt.Errorf("enum[%d] must not be blank", i)
		}
	}

	if len(strings.TrimSpace(s.Default)) == 0 {
		return errors.New("default must not be blank")
	}

	if s.Description != nil && len(strings.TrimSpace(*s.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	return nil
}
