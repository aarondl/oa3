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
			"/path/to/api": &PathRef{
				Path: &Path{
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
	if len(refs) != 6 {
		// In Paths: Pathref (val), RequestBodyRef (val), SchemaRef (ref)
		// In Components: SchemaRef (ref), SchemaRef (ref), SchemaRef (val)
		t.Error("number of refs wrong:", len(refs))
	}
}

func TestResolveRefs(t *testing.T) {
	t.Parallel()

	var testGraph = &OpenAPI3{
		Paths: Paths{
			"/path/to/api": &PathRef{
				Path: &Path{
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
			"/path/to/api": &PathRef{
				Path: &Path{
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
