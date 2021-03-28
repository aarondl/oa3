package tsclient

import (
	"os"
	"testing"

	"github.com/aarondl/fixtures"
	"github.com/aarondl/oa3/openapi3spec"
)

func TestGenerator(t *testing.T) {
	t.Parallel()

	oa, err := openapi3spec.LoadYAML("testdata/petstore.yaml", true)
	if err != nil {
		t.Fatal(err)
	}

	gen := New()

	err = gen.Load(os.DirFS("../templates/tsclient"))
	if err != nil {
		t.Fatal(err)
	}

	files, err := gen.Do(oa, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		fixtures.Bytes(t, f.Name, f.Contents)
	}
}
