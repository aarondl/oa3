// Package support is full of helper functions and types for the
// code generator
package support

import (
	"errors"
	"net/http"
)

var (
	// ErrNoBody is returned from a handler and expected to be handled by
	// ErrorHandler in some useful way for the application.
	ErrNoBody = errors.New("no body")
)

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

// AddErrs adds errors to an error map and returns the map
func AddErrs(errs Errors, key string, toAdd ...error) Errors {
	if len(toAdd) == 0 {
		return errs
	}

	if errs == nil {
		errs = make(map[string][]string)
	}

	fieldErrs := errs[key]
	for _, e := range toAdd {
		fieldErrs = append(fieldErrs, e.Error())
	}
	errs[key] = fieldErrs

	return errs
}

// MergeErrs merges src's keys and values into dst. dst is created if it is nil.
// Returns dst. Colliding keys will be overwritten by what's in src.
func MergeErrs(dst Errors, src Errors) Errors {
	if len(src) == 0 {
		return dst
	}

	if dst == nil {
		dst = make(Errors)
	}

	for k, v := range src {
		errs := make([]string, len(v))
		copy(errs, v)
		dst[k] = errs
	}

	return dst
}
