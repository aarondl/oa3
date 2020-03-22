// Package generator makes a small abstraction for
// main => generators to work with.
package generator

import "github.com/aarondl/oa3/openapi3spec"

// File described by name and contents
type File struct {
	Name     string
	Contents []byte
}

// Interface can load templates and generate file data
type Interface interface {
	Load(templateDir string) error
	Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]File, error)
}
