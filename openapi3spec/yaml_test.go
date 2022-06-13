package openapi3spec

import (
	"testing"

	"github.com/aarondl/fixtures"
)

func TestYAML(t *testing.T) {
	t.Parallel()

	oa, err := LoadYAML("testdata/openapi3.yaml", false)
	if err != nil {
		t.Fatal(err)
	}

	fixtures.JSON(t, "openapi3.yaml", oa)
}

func TestFileRefs(t *testing.T) {
	t.Parallel()

	path := "testdata/multi_file.yaml"

	oa, err := LoadYAML(path, false)
	if err != nil {
		t.Fatal(err)
	}

	if err := oa.ResolveRefs(path); err != nil {
		t.Fatal(err)
	}

	fixtures.JSON(t, "multifile.yaml", oa)
}
