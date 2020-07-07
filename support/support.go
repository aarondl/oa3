// Package support is full of helper functions and types for the
// code generator
package support

import "net/http"

type (
	// Errors is how validation errors are given to the ValidationConverter
	Errors map[string][]string

	// ValidationConverter is used to convert validation errors
	// to something that will work as a response for all methods
	// that must return validation errors.
	ValidationConverter func(Errors) interface{}

	// MW is a middleware stack divided into tags. The first tag of an operation
	// decides what middleware it belongs to. The empty string is middleware
	// for untagged operations.
	MW map[string][]func(w http.ResponseWriter, r *http.Request)
)

// ErrorHandler is an adapter that allows routing to special http.HandlerFuncs
// that additionall have an error return.
type ErrorHandler interface {
	Wrap(func(w http.ResponseWriter, r *http.Request) error) http.Handler
}
