package openapi3spec

import "fmt"

// The Link object represents a possible design-time link for a response. The
// presence of a link does not guarantee the caller's ability to successfully
// invoke it, rather it provides a known relationship and traversal mechanism
// between responses and other operations. Unlike dynamic links (i.e. links
// provided in the response payload), the OAS linking mechanism does not require
// link information in the runtime response. For computing links, and providing
// instructions to execute them, a runtime expression is used for accessing
// values in an operation and using them as parameters while invoking the linked
// operation.
type Link struct {
	// The following two fields are mutually exclusive
	OperationRef *string `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  *string `json:"operationId,omitempty" yaml:"operationId,omitempty"`

	Parameters  map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody interface{}            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`

	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	Server      *Server `json:"server,omitempty" yaml:"server,omitempty"`
}

// Validate a link
func (l *Link) Validate(c Components) error {
	if l.OperationRef != nil && l.OperationID != nil {
		return fmt.Errorf("operationRef is mutually exclusive with operationId")
	}

	//TODO: Finish validation of this odd thing

	return nil
}

// LinkRef refers to a link
type LinkRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Link
}

// Validate link ref
func (l *LinkRef) Validate(c Components) error {
	// Don't validate references
	if l == nil || len(l.Ref) != 0 {
		return nil
	}

	if err := l.Link.Validate(c); err != nil {
		return err
	}

	return nil
}
