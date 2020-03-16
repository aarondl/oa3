package openapi3spec

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	rgxComponentName = regexp.MustCompile(`^[a-zA-Z0-9\.\-_]+$`)
)

// Components specify referenceable reusable components
type Components struct {
	Schemas         map[string]*SchemaRef         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*ResponseRef       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*ParameterRef      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]*ExampleRef        `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBodyRef    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]*HeaderRef         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeRef `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]*LinkRef           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]*CallbackRef       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// ParseRef breaks the uri into relevant parts
func ParseRef(uri, kind string) (name string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("$ref failed to parse(%s): %w", uri, err)
	}

	if len(u.Scheme) != 0 || len(u.Path) != 0 {
		return "", fmt.Errorf("$ref cannot contain non-file-local references: %s", uri)
	}

	if len(u.Fragment) == 0 {
		return "", fmt.Errorf("$ref must contain a url fragment: %s", uri)
	}

	splits := strings.Split(u.Fragment, "/")
	if len(splits) != 3 {
		return "", fmt.Errorf("$ref must contain a url in the form '#/components/TYPE/NAME' but got: %s", uri)
	}

	if splits[0] != "components" {
		return "", fmt.Errorf("$ref must start with '#/components' but got: %s", uri)
	}

	switch splits[1] {
	case "schemas", "responses", "parameters", "examples", "requestBodies",
		"headers", "securitySchemes", "links", "callbacks":
		return splits[2], nil
	}

	return "", fmt.Errorf("$ref can only refer to types: schemas|responses|parameters|examples|requestBodies|headers|securitySchemes|links|callbacks, but got: %s", splits[2])
}

// Validate components
func (c *Components) Validate() error {
	for k, v := range c.Schemas {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("schemas(%s): invalid component key name", k)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("schemas(%s).%w", k, err)
		}
	}
	for k, v := range c.Responses {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("responses(%s): invalid component key name", k)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("responses(%s).%w", k, err)
		}
	}
	for k, v := range c.Parameters {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("parameters(%s): invalid component key name", k)
		}
		if err := v.Validate(nil); err != nil {
			return fmt.Errorf("parameters(%s).%w", k, err)
		}
	}
	for k := range c.Examples {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("examples(%s): invalid component key name", k)
		}
		// if err := v.Validate(); err != nil {
		// 	return fmt.Errorf("examples(%s).%w", k, err)
		// }
	}
	for k, v := range c.RequestBodies {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("requestBodies(%s): invalid component key name", k)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("requestBodies(%s).%w", k, err)
		}
	}
	for k, v := range c.Headers {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("headers(%s): invalid component key name", k)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("headers(%s).%w", k, err)
		}
	}
	for k := range c.SecuritySchemes {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("securitySchemes(%s): invalid component key name", k)
		}
		// if err := v.Validate(); err != nil {
		// 	return fmt.Errorf("securitySchemes(%s).%w", k, err)
		// }
	}
	for k, v := range c.Links {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("links(%s): invalid component key name", k)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("links(%s).%w", k, err)
		}
	}
	for k, v := range c.Callbacks {
		if !rgxComponentName.MatchString(k) {
			return fmt.Errorf("callbacks(%s): invalid component key name", k)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("callbacks(%s).%w", k, err)
		}
	}

	return nil
}
