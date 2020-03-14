package openapi3spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	rgxSemver    = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	rgxEmail     = regexp.MustCompile(` [^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+`)
	rgxPathNames = regexp.MustCompile(`\{([a-z_0-9]+)\}`)
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

// Validate the openapi3 object
//
// Although validate sounds like a read-only operation, it sets default values
// according to the spec, and resolves references.
func (o *OpenAPI3) Validate() error {
	if o.Components == nil {
		o.Components = new(Components)
	}

	if !rgxSemver.MatchString(o.OpenAPI) {
		return errors.New("openapi must be a semantic version number")
	}

	if err := o.Info.Validate(); err != nil {
		return err
	}

	if len(o.Servers) == 0 {
		o.Servers = []Server{
			Server{URL: "/"},
		}
	} else {
		for i, s := range o.Servers {
			if err := s.Validate(); err != nil {
				return fmt.Errorf("servers[%d].%w", i, err)
			}
		}
	}

	opIDs := make(map[string]struct{})

	if len(o.Paths) == 0 {
		return errors.New("must have at least one item in top-level 'paths'")
	}
	for k, p := range o.Paths {
		if !strings.HasPrefix(k, "/") {
			return fmt.Errorf("paths(%s): must begin with /", k)
		}

		nameMatches := rgxPathNames.FindAllStringSubmatch(k, -1)
		var names []string
		if len(nameMatches) != 0 {
			sort.Strings(names)
			for i := 0; i < len(names)-1; i++ {
				if names[i] == names[i+1] {
					return fmt.Errorf("paths(%s): has duplicate path parameter: %s", k, names[i])
				}
			}

			names = make([]string, len(nameMatches))
			for i, n := range nameMatches {
				names[i] = n[1]
			}
		}

		if err := p.Validate(names, opIDs); err != nil {
			return fmt.Errorf("paths(%s).%w", k, err)
		}
	}

	return nil
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

// Validate info struct
func (i *Info) Validate() error {
	if i == nil {
		return nil
	}

	if len(strings.TrimSpace(i.Title)) == 0 {
		return errors.New("info.title must not be blank")
	}
	if i.Description != nil && len(strings.TrimSpace(*i.Description)) == 0 {
		return errors.New("info.description if present must not be blank")
	}
	if i.TermsOfService != nil {
		_, err := url.Parse(*i.TermsOfService)
		if err != nil {
			return fmt.Errorf("info.termsOfService if present must be a url: %w", err)
		}
	}

	if err := i.Contact.Validate(); err != nil {
		return err
	}
	if err := i.License.Validate(); err != nil {
		return err
	}

	if len(strings.TrimSpace(i.Version)) == 0 {
		return errors.New("info.version must not be blank")
	}

	return nil
}

// Extensions are for x- extensions to the open api spec, they can be
// any json value
type Extensions map[string]interface{}
