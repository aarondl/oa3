package templates

import "github.com/aarondl/oa3/openapi3spec"

// TemplateData for all generators
type TemplateData struct {
	Spec   *openapi3spec.OpenAPI3
	Params map[string]string

	Imports map[string]struct{}

	Name   string
	Object interface{}
}

func newData(old TemplateData, name string, obj interface{}) TemplateData {
	copy := old
	copy.Name = name
	copy.Object = obj
	return copy
}

func recurseData(old TemplateData, nextName string, nextObj interface{}) TemplateData {
	copy := old
	copy.Name = old.Name + nextName
	copy.Object = nextObj
	return copy
}

// NewTemplateData constructor
func NewTemplateData(spec *openapi3spec.OpenAPI3, params map[string]string) *TemplateData {
	return &TemplateData{
		Spec:    spec,
		Params:  params,
		Imports: make(map[string]struct{}),
	}
}

// NewTemplateDataWithObject constructor
func NewTemplateDataWithObject(spec *openapi3spec.OpenAPI3, params map[string]string, name string, object interface{}) *TemplateData {
	return &TemplateData{
		Spec:    spec,
		Params:  params,
		Imports: make(map[string]struct{}),
		Name:    name,
		Object:  object,
	}
}

// Import records the importing of a library
func (t TemplateData) Import(importName string) string {
	t.Imports[importName] = struct{}{}
	return ""
}
