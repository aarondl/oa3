package openapi3spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	rgxSemver    = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	rgxEmail     = regexp.MustCompile(`^[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+$`)
	rgxPathNames = regexp.MustCompile(`\{([a-z_0-9]+)\}`)
)

// LoadYAML file
//
// Optionally post-process by validating, resolving references, etc.
func LoadYAML(filename string, postProcess bool) (*OpenAPI3, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	oa, err := LoadYAMLReader(file, postProcess)
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

	oa, err := LoadJSONReader(file, postProcess)
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
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	oa := new(OpenAPI3)
	if err = yaml.Unmarshal(b, &oa); err != nil {
		return nil, err
	}

	if postProcess {
		if err = oa.ResolveRefs(); err != nil {
			return nil, fmt.Errorf("error resolving references: %w", err)
		}
		if err = oa.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		oa.CopyInheritedItems()
	}

	return oa, nil
}

// LoadJSONReader loads the open api spec from a json reader
//
// Optionally post-process by validating, resolving references, etc.
func LoadJSONReader(reader io.Reader, postProcess bool) (*OpenAPI3, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	oa := new(OpenAPI3)
	if err = json.Unmarshal(b, &oa); err != nil {
		return nil, err
	}

	if postProcess {
		if err = oa.ResolveRefs(); err != nil {
			return nil, fmt.Errorf("error resolving references: %w", err)
		}
		if err = oa.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		oa.CopyInheritedItems()
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

// CopyInheritedItems effectively pushes higher-level elements down into
// child elements where specified by the spec.
//
// This should be called after Validate & ResolveRefs
//
//   The server array in a child completely overrides any inherited value.
//   - OpenAPI3.Servers -> OpenAPI3.Paths.Servers
//   - OpenAPI3.Paths.Servers -> OpenAPI3.Paths.Operations.(GET|POST..).Servers
//
//    Each element is considered, a duplicating element in the child overrides
//   the parent.
//   - OpenAPI3.Paths.Parameters -> OpenAPI3.Paths.Operations.(GET|POST).Parameters
func (o *OpenAPI3) CopyInheritedItems() {
	for _, path := range o.Paths {
		if len(o.Servers) > 0 && len(path.Servers) == 0 {
			path.Servers = make([]Server, len(o.Servers))
			copy(path.Servers, o.Servers)
		}

		ops := []*Operation{
			path.Get, path.Put, path.Post, path.Delete, path.Options, path.Head,
			path.Patch, path.Trace,
		}

		for _, op := range ops {
			if op == nil {
				continue
			}

			if len(op.Servers) == 0 {
				op.Servers = make([]Server, len(path.Servers))
				copy(op.Servers, path.Servers)
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

// ResolveRefs finds all the $ref's in the spec and uses the
// components to set them. Currently the only supported references are local
// references to the components object.
//
// In order to resolve the references properly we first find all the references.
// These references become an directed acyclic graph (or it should be, we do
// cycle checking just in case). We can then perform a DFS on the references
// we found to determine proper resolution order.
//
// Once we have the order, we simply iterate and set pointers to the correct
// objects.
func (o *OpenAPI3) ResolveRefs() error {
	refs := findAllRefs(reflect.ValueOf(o))

	resolveOrder := make([]interface{}, 0, len(refs))

	for _, ref := range refs {
		var err error
		resolveOrder, err = o.resolveRefsDFS(ref, resolveOrder)
		if err != nil {
			return err
		}
		debugln()
	}

	for i := len(resolveOrder) - 1; i >= 0; i-- {
		ref := reflect.ValueOf(resolveOrder[i])
		refStructType := ref.Type()
		ref = ref.Elem()
		debugf("resolving %s (%s)\n", ref.Type(), ref.Field(0).Interface().(string))

		if !ref.Field(1).IsNil() {
			// Already resolved
			continue
		}

		refName := ref.Field(0).Interface().(string)

		lookup, err := o.lookupReference(refStructType, refName)
		if err != nil {
			return err
		}

		ref.Field(1).Set(lookup.Elem().Field(1))
	}

	return nil
}

func (o *OpenAPI3) resolveRefsDFS(ref interface{}, order []interface{}) ([]interface{}, error) {
	stack := []interface{}{ref}
	seen := make(map[interface{}]bool)

DFS:
	for len(stack) > 0 {
		ref := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		refValue := reflect.ValueOf(ref)
		debugf("pop:     %s (%s)\n", refValue.Elem().Type(), refValue.Elem().Field(0).Interface().(string))

		// They all must be pointers
		refStructType := reflect.TypeOf(ref)
		refStruct := reflect.ValueOf(ref).Elem()
		refURI := refStruct.Field(0).Interface().(string)

		if seen[ref] {
			return nil, fmt.Errorf("cycle detected: %s", refURI)
		}
		seen[ref] = true

		for _, v := range order {
			if v == ref {
				debugf("done(v): %s (%s)\n", refValue.Elem().Type(), refValue.Elem().Field(0).Interface().(string))
				continue DFS
			}
		}

		order = append(order, ref)
		hasValue := !refStruct.Field(1).IsNil()
		if hasValue {
			debugf("done(h): %s (%s)\n", refValue.Elem().Type(), refValue.Elem().Field(0).Interface().(string))
			continue
		}

		debugf("resolve: %s (%s)\n", refValue.Elem().Type(), refValue.Elem().Field(0).Interface().(string))
		// Here we must actually resolve this
		// Find the pointer to set this reference to

		lookup, err := o.lookupReference(refStructType, refURI)
		if err != nil {
			return nil, err
		}

		debugf("push:    %s (%s)\n", lookup.Elem().Type(), lookup.Elem().Field(0).Interface().(string))
		stack = append(stack, lookup.Interface())
	}

	return order, nil
}

func (o *OpenAPI3) lookupReference(refStructType reflect.Type, refURI string) (reflect.Value, error) {
	uri, err := url.Parse(refURI)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): %w", refURI, err)
	} else if len(uri.Scheme) != 0 || len(uri.Path) != 0 || len(uri.Host) != 0 {
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): only fragment style refs supported", refURI)
	} else if len(uri.Fragment) == 0 {
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): must include fragment", refURI)
	}

	splits := strings.Split(uri.Fragment, "/")
	if len(splits) != 4 || splits[0] != "" || splits[1] != "components" {
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): fragment should be #/components/TYPE/NAME", refURI)
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
	default:
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): fragment did not contain a valid type (%s)", refURI, refType)
	}

	lookup := refMap.MapIndex(reflect.ValueOf(refName))
	if !lookup.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): could not find object %s.%s", refURI, refType, refName)
	}

	if lookup.Type() != refStructType {
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): %T references %s", refURI, refStructType.Name(), refType)
	}

	if lookup.IsNil() {
		return reflect.Value{}, fmt.Errorf("invalid ref(%s): struct referred to is nil", refURI)
	}

	return lookup, nil
}

func findAllRefs(val reflect.Value) []interface{} {
	kind := val.Kind()
	if kind == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		kind = val.Kind()
	}

	refs := make([]interface{}, 0)

	switch kind {
	case reflect.Struct:
		ln := val.NumField()

		typ := val.Type()

		if ln == 2 && typ.Field(0).Name == "Ref" {
			// Reference struct, add it to our list of refs
			refs = append(refs, val.Addr().Interface())
		}

		for i := 0; i < ln; i++ {
			field := val.Field(i)
			//debugf("findref(struct): %T\n", field.Interface())
			fieldRefs := findAllRefs(field)
			refs = append(refs, fieldRefs...)
		}
	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			//debugf("findref(map): %s\n", iter.Key().Interface())
			mapRefs := findAllRefs(iter.Value())
			refs = append(refs, mapRefs...)
		}
	case reflect.Slice:
		ln := val.Len()
		for i := 0; i < ln; i++ {
			//debugf("findref(slice): %d\n", i)
			sliceRefs := findAllRefs(val.Index(i))
			refs = append(refs, sliceRefs...)
		}
	default:
		//debugf("skipping %s\n", val.Type().Name())
	}

	return refs
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
