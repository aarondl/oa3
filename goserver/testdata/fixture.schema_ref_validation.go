// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"strings"

	"github.com/aarondl/oa3/support"
)

// Check for arrays that can call validate function
type RefValidation struct {
	MustValidateItem string `json:"mustValidateItem"`
}

// validateSchema validates the object and returns
// errors that can be returned to the user.
func (o RefValidation) validateSchema() support.Errors {
	var ctx []string
	var ers []error
	var errs support.Errors
	_, _, _ = ctx, ers, errs

	ers = nil
	if err := support.ValidateMaxLength(o.MustValidateItem, 5); err != nil {
		ers = append(ers, err)
	}
	if len(ers) != 0 {
		ctx = append(ctx, "mustValidateItem")
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	return errs
}
