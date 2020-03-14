package openapi3spec

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// LoadYAML file
func LoadYAML(filename string) (*OpenAPI3, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	oa, err := LoadYAMLReader(file)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return oa, nil
}

// LoadJSON file
func LoadJSON(filename string) (*OpenAPI3, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	oa, err := LoadJSONReader(file)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return oa, nil
}

// LoadYAMLReader loads the open api spec from a yaml reader
func LoadYAMLReader(reader io.Reader) (*OpenAPI3, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	oa := new(OpenAPI3)
	if err = yaml.Unmarshal(b, &oa); err != nil {
		return nil, err
	}

	return oa, nil
}

// LoadJSONReader loads the open api spec from a json reader
func LoadJSONReader(reader io.Reader) (*OpenAPI3, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	oa := new(OpenAPI3)
	if err = json.Unmarshal(b, &oa); err != nil {
		return nil, err
	}

	return oa, nil
}

// OpenAPI3 is the root of the OpenAPI Document
type OpenAPI3 struct {
	OpenAPI      string                `json:"openapi,omitempty" yaml:"openapi,omitempty"`
	Info         Info                  `json:"info,omitempty" yaml:"info,omitempty"`
	Servers      []Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        Paths                 `json:"paths,omitempty" yaml:"paths,omitempty"`
	Components   *Components           `json:"components,omitempty" yaml:"components,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []Tag                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Info object provides metadata about the API. The metadata MAY be used by the
// clients if needed, and MAY be presented in editing or documentation
// generation tools for convenience.
type Info struct {
	Title          string   `json:"title,omitempty" yaml:"title,omitempty"`
	Description    *string  `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService *string  `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version,omitempty" yaml:"version,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Extensions are for x- extensions to the open api spec, they can be
// any json value
type Extensions map[string]interface{}
