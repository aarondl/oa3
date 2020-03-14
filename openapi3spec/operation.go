package openapi3spec

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Operation on a url
type Operation struct {
	Tags         []string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      *string       `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  *string       `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	OperationID string                  `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  []ParameterRef          `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBodyRef         `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   Responses               `json:"responses,omitempty" yaml:"responses,omitempty"`
	Callbacks   map[string]*CallbackRef `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`

	Deprecated bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security   []SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Servers    []Server              `json:"servers,omitempty" yaml:"servers,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate an operation
func (o *Operation) Validate(pathTemplates []string, opIDs map[string]struct{}) error {
	if o == nil {
		return nil
	}

	sort.Strings(o.Tags)
	for i := 0; i < len(o.Tags)-1; i++ {
		if o.Tags[i] == o.Tags[i+1] {
			return fmt.Errorf("tags has duplicate: %s", o.Tags[i])
		}
	}

	if o.Summary != nil && len(strings.TrimSpace(*o.Summary)) == 0 {
		return errors.New("summary if present must not be blank")
	}
	if o.Description != nil && len(strings.TrimSpace(*o.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	if len(strings.TrimSpace(o.OperationID)) == 0 {
		return errors.New("operationId must not be blank")
	}

	if _, ok := opIDs[o.OperationID]; ok {
		return fmt.Errorf("operationId is duplicated: %s", o.OperationID)
	}
	opIDs[o.OperationID] = struct{}{}

	if err := paramDuplicateKeyCheck(o.Parameters); err != nil {
		return fmt.Errorf("parameters.%w", err)
	}
	for i, p := range o.Parameters {
		if err := p.Validate(pathTemplates); err != nil {
			return fmt.Errorf("parameters[%d].%w", i, err)
		}
	}

	if err := o.RequestBody.Validate(); err != nil {
		return fmt.Errorf("requestBody.%w", err)
	}

	if len(o.Responses) == 0 {
		return errors.New("responses must not be empty")
	}
	for i, r := range o.Responses {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("responses(%s).%w", i, err)
		}
	}
	for k, c := range o.Callbacks {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("callbacks(%s).%w", k, err)
		}
	}
	for i, s := range o.Security {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("security[%d].%w", i, err)
		}
	}
	for i, s := range o.Servers {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("servers[%d].%w", i, err)
		}
	}

	return nil
}
