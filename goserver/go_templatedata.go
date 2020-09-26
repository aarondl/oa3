package goserver

import (
	"net/http"
	"sort"

	"github.com/aarondl/oa3/openapi3spec"
)

type tagPath struct {
	Tag string
	Ops []tagOp
}

type tagOp struct {
	Path   string
	Method string
	Op     *openapi3spec.Operation
}

func tagPaths(spec *openapi3spec.OpenAPI3) ([]tagPath, error) {
	tags := make(map[string][]tagOp)
	for name, p := range spec.Paths {
		do := func(op *openapi3spec.Operation, method string) {
			if op == nil {
				return
			}
			tag := ""
			if len(op.Tags) != 0 {
				tag = op.Tags[0]
			}

			slice := tags[tag]
			slice = append(slice, tagOp{Path: name, Method: method, Op: op})
			tags[tag] = slice
		}

		do(p.Get, http.MethodGet)
		do(p.Post, http.MethodPost)
		do(p.Put, http.MethodPut)
		do(p.Patch, http.MethodPatch)
		do(p.Delete, http.MethodDelete)
		do(p.Options, http.MethodOptions)
		do(p.Head, http.MethodHead)
		do(p.Trace, http.MethodTrace)
	}

	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tagPathsOut := make([]tagPath, 0, len(keys))
	for _, k := range keys {
		tagOperations := tags[k]
		sort.Slice(tagOperations, func(i, j int) bool {
			return tagOperations[i].Op.OperationID < tagOperations[j].Op.OperationID
		})

		tagPathsOut = append(tagPathsOut, tagPath{Tag: k, Ops: tagOperations})
	}

	return tagPathsOut, nil
}
