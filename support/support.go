// Package support is full of helper functions and types for the
// code generator
package support

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/aarondl/json"
	"github.com/aarondl/oa3/openapi3spec"
)

var (
	// ErrNoBody is returned from a handler and expected to be handled by
	// ErrorHandler in some useful way for the application.
	ErrNoBody = errors.New("no body")

	// buffers is a sync pool of buffers for json marshaling
	buffers = sync.Pool{New: newBuffer}
)

func newBuffer() any {
	return new(bytes.Buffer)
}

// getBuffer retrieves a buffer from the buffer pool
func getBuffer() *bytes.Buffer {
	buf := buffers.Get().(*bytes.Buffer)
	buf.Reset()

	return buf
}

// putBuffer back into the buffer pool
func putBuffer(buf *bytes.Buffer) {
	buffers.Put(buf)
}

type (
	// Errors is how validation errors are given to the ValidationConverter
	Errors map[string][]string

	// ValidationConverter is used to convert validation errors
	// to something that will work as a json response for all methods
	// that must return validation errors. It must implement error to be
	// passed back to the ErrorHandler interface.
	ValidationConverter func(Errors) error

	// MW is a middleware stack divided into tags. The first tag of an operation
	// decides what middleware it belongs to. The empty string is middleware
	// for untagged operations.
	MW map[string][]func(http.Handler) http.Handler

	// Interceptor is a more powerful version of a middleware, they receive
	// parts of the specification to be able to do some more fancy handling.
	//
	// The parsedValues will include (in-order) all the parameters that were
	// parsed out. The body, the headers, query strings, url params
	//
	// In the future values may be more specific to be more friendly.
	Interceptor func(
		w http.ResponseWriter,
		r *http.Request,
		schema *openapi3spec.OpenAPI3,
		path *openapi3spec.Path,
		next func() (any, error),
		respond func(any, error) error,
		values ...any,
	) error
)

// ErrorHandler is an adapter that allows routing to special http.HandlerFuncs
// that additionally have an error return.
type ErrorHandler interface {
	Wrap(func(w http.ResponseWriter, r *http.Request) error) http.Handler
}

// AddErrs adds errors to an error map and returns the map
//
//	eg. {"a": ["1"]}, "a", "2" = {"a": ["1", "2"]}
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

// AddErrsFlatten flattens toAdd by adding key on to the errors inside toAdd
//
//	eg. {"a": ["1"]}, "key", {"b": ["2"]} = {"a": ["1"], "key.b": ["2"]}
func AddErrsFlatten(errs Errors, key string, toAdd Errors) Errors {
	if len(toAdd) == 0 {
		return errs
	}

	if errs == nil {
		errs = make(map[string][]string)
	}

	for field, fieldErrs := range toAdd {
		newKey := key + "." + field
		errs[newKey] = fieldErrs
	}

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

// WriteJSON uses a pool of buffers to write into. This avoids a double
// allocation from using json.Marshal (json.Marshal uses its own internal
// pooled buffer and then copies that to a newly allocated []byte, this way
// we should have pools for both json's internal buffer and our own).
func WriteJSON(w http.ResponseWriter, object any) error {
	buf := getBuffer()
	defer putBuffer(buf)

	marshaller := json.NewEncoder(buf)
	if err := marshaller.Encode(object); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	// Ignore errors that fail to write out to clients because these will
	// generally not be solvable (disconnections etc)
	// We copy -1 because marshaller.Encode produces a newline at the end
	// of each json message
	_, err := io.CopyN(w, buf, int64(buf.Len()-1))
	return err
}

// ReadJSON reads JSON from the body and ensures the body is closed.
// object should be a pointer in order to deserialize properly.
// We copy into a pooled buffer to avoid the allocation from ioutil.ReadAll
// or something similar.
func ReadJSON(r *http.Request, object any) error {
	buf := getBuffer()
	defer putBuffer(buf)

	if _, err := io.Copy(buf, r.Body); err != nil {
		return fmt.Errorf("failed to copy into temp buffer: %w", err)
	}

	if err := r.Body.Close(); err != nil {
		return fmt.Errorf("failed to close body after json read: %w", err)
	}

	if err := json.Unmarshal(buf.Bytes(), object); err != nil {
		return err
	}

	return nil
}

// ReadJSONBuffer is like ReadJSON but does not put back its buffer, instead
// it allows further use, must end with a ReturnJSONBuffer call.
//
// If an error occurs the buffer does not need to be returned.
func ReadJSONBuffer(r *http.Request, object any) (*bytes.Buffer, error) {
	buf := getBuffer()

	if _, err := io.Copy(buf, r.Body); err != nil {
		putBuffer(buf)
		return nil, fmt.Errorf("failed to copy into temp buffer: %w", err)
	}

	if err := r.Body.Close(); err != nil {
		putBuffer(buf)
		return nil, fmt.Errorf("failed to close body after json read: %w", err)
	}

	if err := json.Unmarshal(buf.Bytes(), object); err != nil {
		putBuffer(buf)
		return nil, err
	}

	return buf, nil
}

// ReturnJSONBuffer is called to return a buffer to the pool.
func ReturnJSONBuffer(b *bytes.Buffer) {
	putBuffer(b)
}

// Ptr is useful for taking a normal value and converting to a pointer.
func Ptr[T any](v T) *T {
	return &v
}
