// Package tsclient generates typescript clients for the browser
package tsclient

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/aarondl/oa3/generator"
	"github.com/aarondl/oa3/openapi3spec"
	"github.com/aarondl/oa3/templates"
)

// templates for generation
var tpls = []string{
	"client.tpl",
}

// funcs to use for generation
var funcs = map[string]interface{}{
	"lowerFirst": func(s string) string {
		if len(s) == 0 {
			return ""
		} else if len(s) == 1 {
			return strings.ToLower(s)
		}
		return strings.ToLower(s[0:1]) + s[1:]
	},
}

// generator generates templates for TS
type gen struct {
	tpl *template.Template
}

// New go generator
func New() generator.Interface {
	return &gen{}
}

// Load templates
func (g *gen) Load(dir fs.FS) error {
	var err error
	g.tpl, err = templates.Load(funcs, dir, tpls...)
	return err
}

// Do generation for Typescript.
func (g *gen) Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]generator.File, error) {
	var files []generator.File
	f, err := generateClient(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client: %w", err)
	}
	files = append(files, f...)

	return files, nil
}

func generateClient(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Paths == nil {
		return nil, nil
	}

	files := make([]generator.File, 0)

	apiName := strings.Title(strings.ReplaceAll(spec.Info.Title, " ", ""))
	data := templates.NewTemplateDataWithObject(spec, params, apiName, nil)
	filename := generator.FilenameFromTitle(spec.Info.Title) + ".ts"

	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, "client", data); err != nil {
		return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	// This may be useful for imports later
	fileBytes := new(bytes.Buffer)
	fileBytes.Write(buf.Bytes())

	content := make([]byte, len(fileBytes.Bytes()))
	copy(content, fileBytes.Bytes())
	files = append(files, generator.File{Name: filename, Contents: content})

	return files, nil
}
