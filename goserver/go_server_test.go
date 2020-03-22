package goserver

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/aarondl/fixtures"
	"github.com/aarondl/oa3/openapi3spec"
	"github.com/aarondl/oa3/templates"
)

func TestCamelSnake(t *testing.T) {
	tests := map[string]string{
		"FileName.tpl":          "file_name.tpl",
		"FileIDName.tpl":        "file_id_name.tpl",
		"schema_FileIDName.tpl": "schema_file_id_name.tpl",
		"ID":                    "id",
	}

	for test, want := range tests {
		if got := camelSnake(test); got != want {
			t.Errorf("test: %q, want: %q, got %q", test, want, got)
		}
	}
}

func TestTopSchemas(t *testing.T) {
	t.Parallel()

	oa, err := openapi3spec.LoadYAML("testdata/top_level_schemas.yaml", true)
	if err != nil {
		t.Fatal(err)
	}

	tpl, err := templates.Load(funcs, "../templates/go", tpls...)
	if err != nil {
		t.Fatal(err)
	}

	fileBuffers, err := generateTopLevelSchemas(oa, nil, tpl)
	if err != nil {
		t.Fatal(err)
	}

	all := new(bytes.Buffer)
	for _, f := range fileBuffers {
		fmt.Fprintf(all, "// === %s\n%s\n", f.Name, f.Contents)
	}

	fixtures.Bytes(t, "top_level_schemas.go", all.Bytes())
}
