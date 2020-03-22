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
