package openapi3spec

import (
	"flag"
	"strings"
	"testing"
)

func init() {
	flag.BoolVar(&DebugOutput, "debug", false, "Debug output")
}

var (
	flagUpdateGolden = flag.Bool("golden", false, "Update golden files")
)

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
			"/path/to/ref": &PathRef{
				Ref: "#/components/pathItems/Path",
			},
		},
		Components: &Components{
			Schemas: map[string]*SchemaRef{
				"A": &SchemaRef{Ref: "#/components/schemas/B"},
				"B": &SchemaRef{Ref: "#/components/schemas/C"},
				"C": &SchemaRef{Schema: &Schema{Type: "string"}},
			},
			PathItems: map[string]*PathRef{
				"Path": &PathRef{
					Path: &Path{
						Post: &Operation{
							RequestBody: &RequestBodyRef{
								RequestBody: &RequestBody{
									Content: map[string]*MediaType{
										"application/json": &MediaType{
											Schema: SchemaRef{
												Ref: "#/components/schemas/B",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := testGraph.ResolveRefs(""); err != nil {
		t.Fatal(err)
	}

	if testGraph.Paths["/path/to/api"].Get.RequestBody.Content["application/json"].Schema.Schema == nil {
		t.Error("the path/requestbody schema was not resolved correctly")
	}

	if testGraph.Paths["/path/to/ref"].Post == nil {
		t.Error("the path was not resolved correctly")
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

	if err := testGraph.ResolveRefs(""); err == nil {
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
			"/path/overrider": &PathRef{
				Path: &Path{
					Parameters: []*ParameterRef{
						{Parameter: &Parameter{
							Name:            "X_Parent",
							In:              "header",
							AllowEmptyValue: false,
						}},
					},
					Get: &Operation{
						Parameters: []*ParameterRef{
							// Should override parent
							{Parameter: &Parameter{
								Name:          "X_Parent",
								In:            "header",
								AllowReserved: true,
							}},
						},
					},
				},
			},
			"/path/inheritor": &PathRef{
				Path: &Path{
					// Servers should be inherited
					Parameters: []*ParameterRef{
						{Parameter: &Parameter{
							Name: "X_Parent",
							In:   "header",
						}},
					},
					Get: &Operation{
						Parameters: []*ParameterRef{
							{Parameter: &Parameter{
								Name: "X_Child",
								In:   "header",
							}},
							// X_Parent should be inherited
						},
					},
				},
			},
		},
	}

	testGraph.CopyInheritedItems()

	overrider := testGraph.Paths["/path/overrider"]
	if !overrider.Get.Parameters[0].AllowReserved {
		t.Error("overrider parameter was not marked 'allowReserved', override failed")
	}
}
