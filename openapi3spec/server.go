package openapi3spec

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	rgxServerKeys = regexp.MustCompile(`\{([^\}]+)?\}`)
)

// Server represents a server of this API
type Server struct {
	URL         string                     `json:"url,omitempty" yaml:"url,omitempty"`
	Description *string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate server
func (s *Server) Validate() error {
	if s == nil {
		return nil
	}

	if _, err := url.Parse(s.URL); err != nil {
		return fmt.Errorf("url must be a valid url: %w", err)
	}

	if s.Description != nil && len(strings.TrimSpace(*s.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	urlVariableKeys := FindServerVariablesInURL(s.URL)
	foundVariableObject := false
	for _, urlVariable := range urlVariableKeys {
		for variableObject := range s.Variables {
			if urlVariable == variableObject {
				foundVariableObject = true
				break
			}
		}
		if !foundVariableObject {
			return fmt.Errorf("serverVariable[%s] has no corresponding variables entry", urlVariable)
		}
	}

	for k, v := range s.Variables {
		foundUrlKey := false
		for _, key := range urlVariableKeys {
			if key == k {
				foundUrlKey = true
				break
			}
		}
		if !foundUrlKey {
			return fmt.Errorf("serverVariable[%s] not found in url provided", k)
		}

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

	if s.Enum != nil && len(s.Enum) == 0 {
		return errors.New("enum must not be an empty array")
	}
	for i, e := range s.Enum {
		if len(strings.TrimSpace(e)) == 0 {
			return fmt.Errorf("enum[%d] must not be blank", i)
		}
	}

	if len(strings.TrimSpace(s.Default)) == 0 {
		return errors.New("default is required and must not be blank")
	}
	if len(s.Enum) > 0 {
		foundDefault := false
		for _, e := range s.Enum {
			if e == s.Default {
				foundDefault = true
				break
			}
		}
		if !foundDefault {
			return fmt.Errorf("default value %q is not present in 'enum' property", s.Default)
		}
	}

	if s.Description != nil && len(strings.TrimSpace(*s.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	return nil
}

// FindServerVariablesInURL uses a regular expression to parse out server
// variables of the form {...} in the given string. They may not overlap.
func FindServerVariablesInURL(s string) []string {
	matches := rgxServerKeys.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return nil
	}

	out := make([]string, len(matches))
	for i := range matches {
		out[i] = matches[i][1]
	}

	return out
}
