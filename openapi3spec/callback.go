package openapi3spec

// Callback is a map of possible out-of band callbacks related to the parent
// operation. Each value in the map is a Path Item Object that describes a set
// of requests that may be initiated by the API provider and the expected
// responses. The key value used to identify the callback object is an
// expression, evaluated at runtime, that identifies a URL to use for the
// callback operation.
type Callback map[string]*Path

// CallbackRef refers to a callback
type CallbackRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Callback
}
