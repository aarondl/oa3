package templates

import "github.com/aarondl/oa3/openapi3spec"

// TemplateData for all generators
type TemplateData struct {
	Spec   *openapi3spec.OpenAPI3
	Params map[string]string

	Imports map[string]struct{}
}

// NewTemplateData constructor
func NewTemplateData(spec *openapi3spec.OpenAPI3, params map[string]string) *TemplateData {
	return &TemplateData{
		Spec:    spec,
		Params:  params,
		Imports: make(map[string]struct{}),
	}
}

// Import records the importing of a library
func (t *TemplateData) Import(importName string) {
	t.Imports[importName] = struct{}{}
}
