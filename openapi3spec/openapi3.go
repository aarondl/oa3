package openapi3spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	rgxSemver    = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	rgxEmail     = regexp.MustCompile(`^[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+$`)
	rgxPathNames = regexp.MustCompile(`\{([a-zA-Z_0-9]+)\}`)
)

// LoadYAML file
//
// Optionally post-process by validating, resolving references, etc.
func LoadYAML(filename string, postProcess bool) (*OpenAPI3, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	oa, err := loadYAMLReader(file, postProcess, filename)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return oa, nil
}

// LoadJSON file
//
// Optionally post-process by validating, resolving references, etc.
func LoadJSON(filename string, postProcess bool) (*OpenAPI3, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	oa, err := loadJSONReader(file, postProcess, filename)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return oa, nil
}

// LoadYAMLReader loads the open api spec from a yaml reader
//
// Optionally post-process by validating, resolving references, etc.
func LoadYAMLReader(reader io.Reader, postProcess bool) (*OpenAPI3, error) {
	return loadYAMLReader(reader, postProcess, "")
}

func loadYAMLReader(reader io.Reader, postProcess bool, filename string) (*OpenAPI3, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	oa := new(OpenAPI3)
	if err = yaml.Unmarshal(b, &oa); err != nil {
		return nil, err
	}

	if postProcess {
		if err = oa.PostProcess(filename); err != nil {
			return nil, err
		}
	}

	return oa, nil
}

// LoadJSONReader loads the open api spec from a json reader
//
// Optionally post-process by validating, resolving references, etc.
func LoadJSONReader(reader io.Reader, postProcess bool) (*OpenAPI3, error) {
	return loadJSONReader(reader, postProcess, "")
}

func loadJSONReader(reader io.Reader, postProcess bool, filename string) (*OpenAPI3, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	oa := new(OpenAPI3)
	if err = json.Unmarshal(b, &oa); err != nil {
		return nil, err
	}

	if postProcess {
		if err = oa.PostProcess(filename); err != nil {
			return nil, err
		}
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

// PostProcess the loaded openapi3 document
func (o *OpenAPI3) PostProcess(filename string) error {
	var err error
	if err = o.ResolveRefs(filename); err != nil {
		return fmt.Errorf("error resolving references: %w", err)
	}
	if err = o.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	o.CopyInheritedItems()
	if err = o.ResolveAllOfs(); err != nil {
		return fmt.Errorf("error resolving allOfs: %w", err)
	}

	return nil
}

// CopyInheritedItems effectively pushes higher-level elements down into
// child elements where specified by the spec.
//
// This should be called after Validate & ResolveRefs
//
//    Each element is considered, a duplicating element in the child overrides
//   the parent.
//   - OpenAPI3.Paths.Parameters -> OpenAPI3.Paths.Operations.(GET|POST).Parameters
func (o *OpenAPI3) CopyInheritedItems() {
	for _, path := range o.Paths {
		ops := []*Operation{
			path.Get, path.Put, path.Post, path.Delete, path.Options, path.Head,
			path.Patch, path.Trace,
		}

		for _, op := range ops {
			if op == nil {
				continue
			}

			for _, pathParam := range path.Parameters {

				found := false
				for _, opParam := range op.Parameters {
					if pathParam.Name == opParam.Name && pathParam.In == opParam.In {
						found = true
						break
					}
				}

				// Overridden by the op
				if found {
					continue
				}

				op.Parameters = append(op.Parameters, pathParam)
			}
		}
	}
}

// ResolveRefs finds all the $ref's in the spec and attempts to set them to
// real values.
//
// In order to do so we use a simple recursive DFS and resolve as we go, because
// this must be an acyclic graph we also fail gracefully on cycles.
func (o *OpenAPI3) ResolveRefs(filename string) error {
	return o.resolveRefsDFS(o, 0, make(map[string]bool), filename)
}

func (o *OpenAPI3) resolveRefsDFS(node any, depth int, processing map[string]bool, filename string) error {
	nodeVal := reflect.ValueOf(node)
	nodeType := nodeVal.Type()
	debugf("%svisit:     %s\n", strings.Repeat(" ", 2*depth), nodeVal.Type())
	if nodeVal.Kind() == reflect.Ptr {
		nodeVal = nodeVal.Elem()
	}

	kind := nodeVal.Kind()
	switch kind {
	case reflect.Struct:
		ln := nodeVal.NumField()
		for i := 0; i < ln; i++ {
			field := nodeVal.Field(i)
			debugf("%sdfs(struct): %T\n", strings.Repeat(" ", 2*depth), field.Interface())
			if field.Kind() != reflect.Ptr {
				field = field.Addr()
			}
			if err := o.resolveRefsDFS(field.Interface(), depth+1, processing, filename); err != nil {
				return err
			}
		}

		// Check if we are a reference struct, in which case we must resolve
		// it.
		if ln == 2 && nodeType.Elem().Field(0).Name == "Ref" {
			refURI := nodeVal.Field(0).Interface().(string)
			debugf("%sref(%s): %s %s\n", strings.Repeat(" ", 2*depth), refURI, nodeType, nodeType.Elem().Field(1).Type)

			// If the ref already has a value that means it's already
			// been completed or rather that it does not require resolution
			// it may require recursive resolution for refs inside of itself
			// but that will have been handled by the loop above.
			hasValue := !nodeVal.Field(1).IsNil()
			if hasValue {
				debugf("%sdone: %s (%s)\n", strings.Repeat(" ", 2*depth), nodeVal.Type(), refURI)
				return nil
			}

			// When we begin processing a reference uri, we store it to check
			// for cycles while we process.
			if processing[refURI] {
				return fmt.Errorf("cycle detected: %s depth %d", refURI, depth)
			}
			processing[refURI] = true

			// Fetch a ref pointer value from the referenced location
			debugf("%sresolve: %s (%s)\n", strings.Repeat(" ", 2*depth), nodeVal.Type(), nodeVal.Field(0).Interface().(string))
			lookup, newFilename, err := o.lookupReference(nodeType, refURI, filename)
			if err != nil {
				return err
			}

			// Before we can resolve ourselves, we have to ensure that the
			// refs inside the object we've resolved are also resolved.
			if err := o.resolveRefsDFS(lookup.Interface(), depth+1, processing, newFilename); err != nil {
				return err
			}

			// After all the refs are resolved inside the lookup set our current
			// ref to the completely resolved value.
			nodeVal.Field(1).Set(lookup.Elem().Field(1))

			delete(processing, refURI)
		}
	case reflect.Map:
		iter := nodeVal.MapRange()
		for iter.Next() {
			// debugf("dfs(map): %s\n", strings.Repeat(" ", 2*depth), iter.Key().Interface())
			if err := o.resolveRefsDFS(iter.Value().Interface(), depth+1, processing, filename); err != nil {
				return err
			}
		}
	case reflect.Slice:
		ln := nodeVal.Len()
		for i := 0; i < ln; i++ {
			// debugf("dfs(slice): %d\n", strings.Repeat(" ", 2*depth), i)
			val := nodeVal.Index(i)
			if val.Kind() != reflect.Ptr {
				val = val.Addr()
			}
			if err := o.resolveRefsDFS(val.Interface(), depth+1, processing, filename); err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *OpenAPI3) lookupReference(refStructType reflect.Type, refURI string, filename string) (reflect.Value, string, error) {
	uri, err := url.Parse(refURI)
	if err != nil {
		return reflect.Value{}, "", fmt.Errorf("invalid ref(%s): %w", refURI, err)
	}

	if err := validateRefURI(refURI, uri); err != nil {
		return reflect.Value{}, "", err
	}

	if len(uri.Fragment) == 0 {
		return o.lookupFileRef(refStructType, uri, refURI, filename)
	}

	splits := strings.Split(uri.Fragment, "/")
	if len(splits) != 4 || splits[0] != "" || splits[1] != "components" {
		return reflect.Value{}, "", fmt.Errorf("invalid ref(%s): fragment should be #/components/TYPE/NAME", refURI)
	}

	refType := splits[2]
	refName := splits[3]

	var refMap reflect.Value

	switch refType {
	case "schemas":
		refMap = reflect.ValueOf(o.Components.Schemas)
	case "responses":
		refMap = reflect.ValueOf(o.Components.Responses)
	case "parameters":
		refMap = reflect.ValueOf(o.Components.Parameters)
	case "examples":
		refMap = reflect.ValueOf(o.Components.Examples)
	case "requestBodies":
		refMap = reflect.ValueOf(o.Components.RequestBodies)
	case "headers":
		refMap = reflect.ValueOf(o.Components.Headers)
	case "securitySchemes":
		refMap = reflect.ValueOf(o.Components.SecuritySchemes)
	case "links":
		refMap = reflect.ValueOf(o.Components.Links)
	case "callbacks":
		refMap = reflect.ValueOf(o.Components.Callbacks)
	case "pathItems":
		refMap = reflect.ValueOf(o.Components.PathItems)
	default:
		return reflect.Value{}, "", fmt.Errorf("invalid ref(%s): fragment did not contain a valid type (%s)", refURI, refType)
	}

	lookup := refMap.MapIndex(reflect.ValueOf(refName))
	if !lookup.IsValid() {
		return reflect.Value{}, "", fmt.Errorf("invalid ref(%s): could not find object %s.%s", refURI, refType, refName)
	}

	if lookup.Type() != refStructType {
		return reflect.Value{}, "", fmt.Errorf("invalid ref(%s): %s references %s (%s)", refURI, refStructType, lookup.Type(), refType)
	}

	if lookup.IsNil() {
		return reflect.Value{}, "", fmt.Errorf("invalid ref(%s): struct referred to is nil", refURI)
	}

	return lookup, "", nil
}

func (o *OpenAPI3) lookupFileRef(refStructType reflect.Type, uri *url.URL, refURI string, filename string) (reflect.Value, string, error) {
	newRef := reflect.New(refStructType.Elem())
	newVal := reflect.New(refStructType.Elem().Field(1).Type.Elem())

	var path = uri.Path
	var err error
	if !filepath.IsAbs(uri.Path) {
		var base = filename
		if len(base) == 0 {
			base, err = os.Getwd()
			if err != nil {
				return reflect.Value{}, "", fmt.Errorf("error resolving ref(%s): determine working directory: %w", refURI, err)
			}
		} else {
			base = filepath.Dir(filename)
		}

		path = filepath.Join(base, path)
		if err != nil {
			return reflect.Value{}, "", fmt.Errorf("error resolving ref(%s): could not compute final path using base: %s and path %s", refURI, base, path)
		}
	}

	switch {
	case strings.HasSuffix(path, ".json"):
		return reflect.Value{}, "", errors.New("json loading not implemented")
	case strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"):
		b, err := os.ReadFile(path)
		if err != nil {
			return reflect.Value{}, "", fmt.Errorf("error resolving ref(%s): failed to resolve yaml file ref, could not read file: %w", refURI, err)
		}

		var untyped map[string]any
		if err := yaml.Unmarshal(b, &untyped); err != nil {
			return reflect.Value{}, "", fmt.Errorf("error resolving ref(%s): failed to resolve yaml file ref, could not unmarshal yaml: %w", refURI, err)
		}

		if err := yamlStruct(newVal, untyped); err != nil {
			return reflect.Value{}, "", fmt.Errorf("error resolving ref(%s): failed to merge unmarshalled data into struct: %w", refURI, err)
		}
	default:
		return reflect.Value{}, "", fmt.Errorf("invalid ref(%s): only yaml/json file refs are supported (must have proper file extension)", refURI)
	}

	newRef.Elem().Field(1).Set(newVal)
	return newRef, path, nil
}

func validateRefURI(refURI string, uri *url.URL) error {
	if (len(uri.Scheme) == 0 || uri.Scheme == "file://") && len(uri.Fragment) == 0 && (len(uri.Host) == 0 || uri.Host == "localhost") && len(uri.Path) != 0 {
		return nil
	}
	if len(uri.Fragment) != 0 && len(uri.Scheme) == 0 && (len(uri.Host) == 0 || uri.Host == "localhost") && len(uri.Path) == 0 {
		return nil
	}

	return fmt.Errorf("invalid ref(%s): only #/fragment or (file://localhost)?/path/to/file.(yaml|yml|json) refs supported", refURI)
}

// Validate the openapi3 object
//
// Although validate sounds like a read-only operation, it also sets default
// values according to the spec.
//
// It should be called after references are resolved.
func (o *OpenAPI3) Validate() error {
	if !rgxSemver.MatchString(o.OpenAPI) {
		return errors.New("openapi must be a semantic version number")
	}

	if !strings.HasPrefix(o.OpenAPI, "3.") && !strings.HasPrefix(o.OpenAPI, "v3.") {
		return errors.New("openapi version must be 3.x.x for use with this package")
	}

	if err := o.Info.Validate(); err != nil {
		return err
	}

	if len(o.Servers) == 0 {
		o.Servers = []Server{
			{URL: "/"},
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

	if err := o.Components.Validate(); err != nil {
		return fmt.Errorf("components.%w", err)
	}

	for i, s := range o.Security {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("security[%d].%w", i, err)
		}
	}

	for i, t := range o.Tags {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("tags[%d].%w", i, err)
		}
	}

	if err := o.ExternalDocs.Validate(); err != nil {
		return fmt.Errorf("externalDocs.%w", err)
	}

	return nil
}

// ResolveAllOfs combines any object that declares it is an 'allOf' into a
// single object.
func (o *OpenAPI3) ResolveAllOfs() error {
	mergeSchema := func(schema *Schema) error {
		if schema == nil {
			return nil
		}

		if schema.AllOf == nil {
			return nil
		}

		schema.Required = make([]string, 0, 0)
		schema.Properties = make(map[string]*SchemaRef)
		for _, s := range schema.AllOf {
			schema.Required = append(schema.Required, s.Schema.Required...)
			for k, v := range s.Properties {
				schema.Properties[k] = v
			}
			if s.AdditionalProperties != nil {
				schema.AdditionalProperties = s.AdditionalProperties
			}
			if s.Discriminator != nil {
				schema.Discriminator = s.Discriminator
			}
		}

		return nil
	}

	checkMediaTypes := func(medias map[string]*MediaType) error {
		for media, m := range medias {
			if err := mergeSchema(m.Schema.Schema); err != nil {
				return fmt.Errorf("content(%s).%w", media, mergeSchema(m.Schema.Schema))
			}
		}

		return nil
	}

	if o.Components != nil {
		for k, v := range o.Components.Schemas {
			if err := mergeSchema(v.Schema); err != nil {
				return fmt.Errorf("components.schemas(%s).%w", k, err)
			}
		}
		for k, v := range o.Components.RequestBodies {
			if v.RequestBody != nil {
				if err := checkMediaTypes(v.RequestBody.Content); err != nil {
					return fmt.Errorf("components.requestBodies(%s).%w", k, err)
				}
			}
		}
		for k, v := range o.Components.Responses {
			if v.Response != nil {
				if err := checkMediaTypes(v.Response.Content); err != nil {
					return fmt.Errorf("components.responses(%s).%w", k, err)
				}
			}
		}
		for k, v := range o.Components.Parameters {
			if err := mergeSchema(v.Schema.Schema); err != nil {
				return fmt.Errorf("components.parameters(%s).%w", k, err)
			}
			if err := checkMediaTypes(v.Content); err != nil {
				return fmt.Errorf("components.parameters(%s).%w", k, err)
			}
		}
	}

	for k, v := range o.Paths {
		ops := map[string]*Operation{
			"get":     v.Get,
			"post":    v.Post,
			"put":     v.Put,
			"delete":  v.Delete,
			"patch":   v.Patch,
			"head":    v.Head,
			"options": v.Options,
		}

		for verb, o := range ops {
			if o == nil {
				continue
			}

			for i, p := range o.Parameters {
				if err := mergeSchema(p.Schema.Schema); err != nil {
					return fmt.Errorf("paths(%s).%s.parameters[%d].%w", k, verb, i, err)
				}
				if err := checkMediaTypes(p.Content); err != nil {
					return fmt.Errorf("paths(%s).%s.parameters[%d].%w", k, verb, i, err)
				}
			}
			if o.RequestBody != nil {
				if err := checkMediaTypes(o.RequestBody.Content); err != nil {
					return fmt.Errorf("paths(%s).%s.requestBody.%w", k, verb, err)
				}
			}
			for resp, r := range o.Responses {
				if err := checkMediaTypes(r.Content); err != nil {
					return fmt.Errorf("paths(%s).%s.responses(%s).%w", k, verb, resp, err)
				}
			}
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
type Extensions map[string]any
