package templates

import "github.com/aarondl/oa3/openapi3spec"

// TemplateData for all generators
type TemplateData struct {
	Spec   *openapi3spec.OpenAPI3
	Params map[string]string

	Imports map[string]struct{}

	Name     string
	Object   any
	Required bool
}

func newData(old TemplateData, name string, obj any) TemplateData {
	copy := old
	copy.Name = name
	copy.Object = obj
	return copy
}

func newDataRequired(old TemplateData, name string, obj any, required bool) TemplateData {
	copy := old
	copy.Name = name
	copy.Object = obj
	copy.Required = required
	return copy
}

func recurseData(old TemplateData, nextName string, nextObj any) TemplateData {
	copy := old
	copy.Name = old.Name + nextName
	copy.Object = nextObj
	return copy
}

func recurseDataSetRequired(old TemplateData, nextName string, nextObj any, required bool) TemplateData {
	copy := old
	copy.Name = old.Name + nextName
	copy.Object = nextObj
	copy.Required = required
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
func NewTemplateDataWithObject(spec *openapi3spec.OpenAPI3, params map[string]string, name string, object any, required bool) *TemplateData {
	return &TemplateData{
		Spec:     spec,
		Params:   params,
		Imports:  make(map[string]struct{}),
		Name:     name,
		Object:   object,
		Required: required,
	}
}

func templateParamExists(td TemplateData, param string) bool {
	_, ok := td.Params[param]
	return ok
}

func templateParamEquals(td TemplateData, param, want string) bool {
	if val, ok := td.Params[param]; ok && val == param {
		return true
	}
	return false
}

// Import records the importing of a library
func (t TemplateData) Import(importName string) string {
	t.Imports[importName] = struct{}{}
	return ""
}
