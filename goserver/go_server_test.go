package goserver

import (
	"os"
	"testing"

	"github.com/aarondl/fixtures"
	"github.com/aarondl/oa3/openapi3spec"
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

func TestGenerator(t *testing.T) {
	t.Parallel()

	oa, err := openapi3spec.LoadYAML("testdata/go_server.yaml", true)
	if err != nil {
		t.Fatal(err)
	}

	gen := New()

	err = gen.Load(os.DirFS("../templates/go"))
	if err != nil {
		t.Fatal(err)
	}

	files, err := gen.Do(oa, map[string]string{
		"timetype":    "chrono",
		"decimaltype": "shopspring",
		"uuidtype":    "google",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		fixtures.Bytes(t, f.Name, f.Contents)
	}
}

func TestGoCodeGenerator(t *testing.T) {
	t.Parallel()

	oa, err := openapi3spec.LoadYAML("testdata/go_server.yaml", true)
	if err != nil {
		t.Fatal(err)
	}

	str := codeForValue(oa)
	fixtures.String(t, t.Name(), str)
}
