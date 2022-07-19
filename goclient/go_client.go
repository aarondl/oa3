package goclient

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"strings"
	"text/template"

	"github.com/aarondl/oa3/generator"
	"github.com/aarondl/oa3/goserver"
	"github.com/aarondl/oa3/openapi3spec"
	"github.com/aarondl/oa3/templates"
)

// templates for generation
var tpls = []string{
	"client_interface.tpl",
	"client_methods.tpl",
	"schema.tpl",
	"schema_top.tpl",

	"validate_schema_noop.tpl",
}

// generator generates templates for Go
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
	g.tpl, err = templates.Load(goserver.TemplateFunctions, dir, tpls...)
	return err
}

// Do generation for Go.
func (g *gen) Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]generator.File, error) {
	if params == nil {
		params = make(map[string]string)
	}
	if pkg, ok := params[goserver.ParamKeyPackage]; !ok || len(pkg) == 0 {
		params[goserver.ParamKeyPackage] = goserver.DefaultPackage
	}

	var files []generator.File
	f, err := goserver.GenerateTopLevelSchemas(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schemas: %w", err)
	}

	files = append(files, f...)

	f, err = generateClientInterface(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to client interface: %w", err)
	}

	files = append(files, f...)

	for i, f := range files {
		formatted, err := format.Source(f.Contents)
		if err != nil {
			return nil, fmt.Errorf("failed to format file(%s): %w\n%s", f.Name, err, f.Contents)
		}

		files[i].Contents = formatted
	}

	return files, nil
}

func generateClientInterface(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Paths == nil {
		return nil, nil
	}

	files := make([]generator.File, 0)

	apiName := strings.Title(strings.ReplaceAll(spec.Info.Title, " ", ""))

	data := templates.NewTemplateDataWithObject(spec, params, apiName, nil, false)

	filename := generator.FilenameFromTitle(spec.Info.Title) + ".go"

	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, "client_interface", data); err != nil {
		return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	fileBytes := new(bytes.Buffer)
	pkg := params["package"]

	fileBytes.WriteString(goserver.Disclaimer)
	fmt.Fprintf(fileBytes, "\npackage %s\n", pkg)
	if imps := goserver.Imports(data.Imports); len(imps) != 0 {
		fileBytes.WriteByte('\n')
		fileBytes.WriteString(goserver.Imports(data.Imports))
		fileBytes.WriteByte('\n')
	}
	fileBytes.WriteByte('\n')
	fileBytes.Write(buf.Bytes())

	content := make([]byte, len(fileBytes.Bytes()))
	copy(content, fileBytes.Bytes())
	files = append(files, generator.File{Name: filename, Contents: content})

	data = templates.NewTemplateDataWithObject(spec, params, apiName, nil, false)
	filename = generator.FilenameFromTitle(spec.Info.Title) + "_methods.go"

	buf.Reset()
	fileBytes.Reset()
	if err := tpl.ExecuteTemplate(buf, "client_methods", data); err != nil {
		return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	fileBytes.WriteString(goserver.Disclaimer)
	fmt.Fprintf(fileBytes, "\npackage %s\n", pkg)
	if imps := goserver.Imports(data.Imports); len(imps) != 0 {
		fileBytes.WriteByte('\n')
		fileBytes.WriteString(goserver.Imports(data.Imports))
		fileBytes.WriteByte('\n')
	}
	fileBytes.WriteByte('\n')
	fileBytes.Write(buf.Bytes())

	files = append(files, generator.File{Name: filename, Contents: fileBytes.Bytes()})

	return files, nil
}
