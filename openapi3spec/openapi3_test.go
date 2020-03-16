package openapi3spec

import (
	"flag"
	"reflect"
	"strings"
	"testing"
)

func init() {
	flag.BoolVar(&DebugOutput, "debug", false, "Debug output")
}

var (
	flagUpdateGolden = flag.Bool("golden", false, "Update golden files")
)

func TestFindRefs(t *testing.T) {
	t.Parallel()

	var testGraph = &OpenAPI3{
		Paths: Paths{
			"/path/to/api": &Path{
				Get: &Operation{
					RequestBody: &RequestBodyRef{
						RequestBody: &RequestBody{
							Content: map[string]*MediaType{
								"application/json": &MediaType{
									Schema: SchemaRef{
										Ref: "#/components/schemas/A",
									},
								},
							},
						},
					},
				},
			},
		},
		Components: &Components{
			Schemas: map[string]*SchemaRef{
				"A": &SchemaRef{Ref: "#/components/schemas/B"},
				"B": &SchemaRef{Ref: "#/components/schemas/C"},
				"C": &SchemaRef{Schema: &Schema{Type: "string"}},
			},
		},
	}

	refs := findAllRefs(reflect.ValueOf(testGraph))
	if len(refs) != 5 {
		// In Paths: RequestBodyRef (val), SchemaRef (ref)
		// In Components: SchemaRef (ref), SchemaRef (ref), SchemaRef (val)
		t.Error("number of refs wrong:", len(refs))
	}
}

func TestResolveRefs(t *testing.T) {
	t.Parallel()

	var testGraph = &OpenAPI3{
		Paths: Paths{
			"/path/to/api": &Path{
				Get: &Operation{
					RequestBody: &RequestBodyRef{
						RequestBody: &RequestBody{
							Content: map[string]*MediaType{
								"application/json": &MediaType{
									Schema: SchemaRef{
										Ref: "#/components/schemas/A",
									},
								},
							},
						},
					},
				},
			},
		},
		Components: &Components{
			Schemas: map[string]*SchemaRef{
				"A": &SchemaRef{Ref: "#/components/schemas/B"},
				"B": &SchemaRef{Ref: "#/components/schemas/C"},
				"C": &SchemaRef{Schema: &Schema{Type: "string"}},
			},
		},
	}

	if err := testGraph.ResolveRefs(); err != nil {
		t.Fatal(err)
	}

	if testGraph.Paths["/path/to/api"].Get.RequestBody.Content["application/json"].Schema.Schema == nil {
		t.Error("the path/requestbody schema was not resolved correctly")
	}

	if testGraph.Components.Schemas["A"].Schema == nil {
		t.Error("A was not resolved correctly")
	}
	if testGraph.Components.Schemas["B"].Schema == nil {
		t.Error("B was not resolved correctly")
	}
	if testGraph.Components.Schemas["C"].Schema == nil {
		t.Error("C is somehow nil now")
	}
}

func TestCycle(t *testing.T) {
	t.Parallel()

	var testGraph = &OpenAPI3{
		Paths: Paths{
			"/path/to/api": &Path{
				Get: &Operation{
					RequestBody: &RequestBodyRef{
						RequestBody: &RequestBody{
							Content: map[string]*MediaType{
								"application/json": &MediaType{
									Schema: SchemaRef{
										Ref: "#/components/schemas/A",
									},
								},
							},
						},
					},
				},
			},
		},
		Components: &Components{
			Schemas: map[string]*SchemaRef{
				"A": &SchemaRef{Ref: "#/components/schemas/B"},
				"B": &SchemaRef{Ref: "#/components/schemas/A"},
				"C": &SchemaRef{Schema: &Schema{Type: "string"}},
			},
		},
	}

	if err := testGraph.ResolveRefs(); err == nil {
		t.Fatal("it should have given us a cycle error")
	} else if !strings.Contains(err.Error(), "cycle detected") {
		t.Fatal("it should have given us a cycle error, but got:", err)
	}
}

func TestCopyInherited(t *testing.T) {
	t.Parallel()

	var testGraph = &OpenAPI3{
		Servers: []Server{{URL: "a"}},
		Paths: Paths{
			"/path/overrider": &Path{
				Parameters: []*ParameterRef{
					&ParameterRef{Parameter: &Parameter{
						Name:            "X_Parent",
						In:              "header",
						AllowEmptyValue: false,
					}},
				},
				// Should override parent
				Servers: []Server{{URL: "c"}},
				Get: &Operation{
					// Should override parent
					Servers: []Server{{URL: "d"}},

					Parameters: []*ParameterRef{
						// Should override parent
						&ParameterRef{Parameter: &Parameter{
							Name:          "X_Parent",
							In:            "header",
							AllowReserved: true,
						}},
					},
				},
			},
			"/path/inheritor": &Path{
				// Servers should be inherited
				Parameters: []*ParameterRef{
					&ParameterRef{Parameter: &Parameter{
						Name: "X_Parent",
						In:   "header",
					}},
				},
				Get: &Operation{
					Parameters: []*ParameterRef{
						&ParameterRef{Parameter: &Parameter{
							Name: "X_Child",
							In:   "header",
						}},
						// X_Parent should be inherited
					},
				},
			},
		},
	}

	testGraph.CopyInheritedItems()

	overrider := testGraph.Paths["/path/overrider"]
	if overrider.Get.Servers[0].URL != "d" {
		t.Error("overrider server was wrong:", overrider.Get.Servers[0].URL)
	}
	if !overrider.Get.Parameters[0].AllowReserved {
		t.Error("overrider parameter was not marked 'allowReserved', override failed")
	}

	inheritor := testGraph.Paths["/path/inheritor"]
	if inheritor.Servers[0].URL != "a" {
		t.Error("servers should be inherited by the path item:", inheritor.Servers[0].URL)
	}
	if inheritor.Get.Servers[0].URL != "a" {
		t.Error("servers should be inherited by the operation:", inheritor.Get.Servers[0].URL)
	}
}
