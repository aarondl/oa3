package openapi3spec

import "fmt"

// Callback is a map of possible out-of band callbacks related to the parent
// operation. Each value in the map is a Path Item Object that describes a set
// of requests that may be initiated by the API provider and the expected
// responses. The key value used to identify the callback object is an
// expression, evaluated at runtime, that identifies a URL to use for the
// callback operation.
type Callback map[string]*Path

// Validate callback
func (c *Callback) Validate(cmps Components) error {
	for k, p := range *c {
		if err := p.Validate(cmps, nil, nil); err != nil {
			return fmt.Errorf("callback(%s).%w", k, err)
		}
	}

	return nil
}

// CallbackRef refers to a callback
type CallbackRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Callback
}

// Validate response ref
func (c *CallbackRef) Validate(cmps Components) error {
	// Don't validate references
	if c == nil || len(c.Ref) != 0 {
		return nil
	}

	if err := c.Callback.Validate(cmps); err != nil {
		return err
	}

	return nil
}
