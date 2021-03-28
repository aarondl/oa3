package elm

import (
	"io/fs"

	"github.com/aarondl/oa3/generator"
	"github.com/aarondl/oa3/openapi3spec"
)

type gen struct{}

// New generator
func New() generator.Interface {
	return &gen{}
}

// Load templates
func (g *gen) Load(dir fs.FS) error { return nil }

// Do generation for Elm
func (g *gen) Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]generator.File, error) {
	return nil, nil
}
