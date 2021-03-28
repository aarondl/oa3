// Package generator makes a small abstraction for
// main => generators to work with.
package generator

import (
	"io/fs"
	"strings"

	"github.com/aarondl/oa3/openapi3spec"
)

// File described by name and contents
type File struct {
	Name     string
	Contents []byte
}

// Interface can load templates and generate file data
type Interface interface {
	Load(fs fs.FS) error
	Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]File, error)
}

var filenameReplacer = strings.NewReplacer(
	" ", "_",
	"\t", "_",
	"\n", "_",
)

// FilenameFromTitle creates a filename from a title
func FilenameFromTitle(title string) string {
	return strings.ToLower(filenameReplacer.Replace(title))
}
