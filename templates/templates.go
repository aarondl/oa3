package templates

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/aarondl/oa3/openapi3spec"
)

// GlobalFunctions for templates
var GlobalFunctions = map[string]interface{}{
	"refName":                refName,
	"mustValidate":           mustValidate,
	"mustValidateRecurse":    mustValidateRecurse,
	"keysReflect":            keysReflect,
	"httpStatus":             http.StatusText,
	"newData":                newData,
	"newDataRequired":        newDataRequired,
	"recurseData":            recurseData,
	"recurseDataSetRequired": recurseDataSetRequired,
}

func refName(ref string) string {
	splits := strings.Split(ref, "/")
	return splits[len(splits)-1]
}

func keysReflect(mp interface{}) ([]string, error) {
	mapType := reflect.TypeOf(mp)
	if mapType.Kind() != reflect.Map || mapType.Key().Kind() != reflect.String {
		return nil, fmt.Errorf("want map[string]X, got: %s", mapType.Name())
	}

	val := reflect.ValueOf(mp)
	iter := val.MapRange()

	keys := make([]string, 0, val.Len())
	for iter.Next() {
		keys = append(keys, iter.Key().String())
	}

	return keys, nil
}

// mustValidate checks to see if the schema requires any kind of validation
func mustValidate(s *openapi3spec.Schema) bool {
	return s.MultipleOf != nil ||
		s.Maximum != nil ||
		s.Minimum != nil ||
		s.MaxLength != nil ||
		s.MinLength != nil ||
		s.Pattern != nil ||
		s.MaxItems != nil ||
		s.MinItems != nil ||
		s.UniqueItems != nil ||
		s.MaxProperties != nil ||
		s.MinProperties != nil ||
		s.Format != nil
}

// mustValidateRecure checks to see if the current schema, or any sub-schema
// requires validation
func mustValidateRecurse(s *openapi3spec.Schema) bool {
	cycleMarkers := make(map[string]struct{})
	return mustValidateRecurseHelper(s, cycleMarkers)
}

func mustValidateRecurseHelper(s *openapi3spec.Schema, visited map[string]struct{}) bool {
	if mustValidate(s) {
		return true
	}

	if s.Type == "array" {
		if len(s.Items.Ref) != 0 {
			if _, ok := visited[s.Items.Ref]; ok {
				return false
			}
			visited[s.Items.Ref] = struct{}{}
		}

		return mustValidateRecurseHelper(s.Items.Schema, visited)
	} else if s.Type == "object" {
		mustV := false
		if s.AdditionalProperties != nil {
			if len(s.AdditionalProperties.Ref) != 0 {
				if _, ok := visited[s.AdditionalProperties.Ref]; !ok {
					visited[s.AdditionalProperties.Ref] = struct{}{}
					mustV = mustV || mustValidateRecurseHelper(s.AdditionalProperties.Schema, visited)
				}
			} else {
				mustV = mustV || mustValidateRecurseHelper(s.AdditionalProperties.Schema, visited)
			}
		}

		for _, v := range s.Properties {
			if len(v.Ref) != 0 {
				if _, ok := visited[v.Ref]; ok {
					continue
				}
				visited[v.Ref] = struct{}{}
			}

			mustV = mustV || mustValidateRecurseHelper(v.Schema, visited)
		}

		return mustV
	}

	return false
}

// Load takes in funcs to apply to each template, a directory that
// contains the files, and the file paths relative that directory
func Load(generatorFuncs map[string]interface{}, dir fs.FS, files ...string) (*template.Template, error) {
	tpl := template.New("").Funcs(sprig.TxtFuncMap()).Funcs(GlobalFunctions)
	if generatorFuncs != nil {
		tpl = tpl.Funcs(generatorFuncs)
	}

	for _, f := range files {
		name := f

		if base := filepath.Base(name); strings.IndexByte(base, '.') > 0 {
			dot := strings.LastIndexByte(name, '.')
			name = name[:dot]
		}

		file, err := dir.Open(f)
		if err != nil {
			return nil, fmt.Errorf("failed to open template file: %q: %w", f, err)
		}

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read template file %q: %w", f, err)
		}

		if err = file.Close(); err != nil {
			return nil, fmt.Errorf("failed to close template file %q: %w", f, err)
		}

		_, err = tpl.New(name).Parse(string(b))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template file %q: %w", f, err)
		}
	}

	return tpl, nil
}
