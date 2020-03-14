package openapi3spec

import "errors"

// UnmarshalJSON completely overrides the typical recursive json decoder
// behavior with its own ideas about how to unmarshal in order to handle
// some idiosynchracies in the spec.
func (o *OpenAPI3) UnmarshalJSON(in []byte) error {
	return errors.New("not implemented")
}
