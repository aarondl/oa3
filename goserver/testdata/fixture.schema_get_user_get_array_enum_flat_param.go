// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"fmt"
	"strings"

	"github.com/aarondl/oa3/support"
)

type GetUserGetArrayEnumFlatParam []GetUserGetArrayEnumFlatParamItem

// validateSchema validates the object and returns
// errors that can be returned to the user.
func (o GetUserGetArrayEnumFlatParam) validateSchema() support.Errors {
	var ctx []string
	var ers []error
	var errs support.Errors
	_, _, _ = ctx, ers, errs

	for i, o := range o {
		_ = o
		ctx = append(ctx, fmt.Sprintf("[%d]", i))
		var ers []error

		ers = nil
		if err := support.ValidateEnum(o, []string{"a", "b"}); err != nil {
			ers = append(ers, err)
		}

		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	return errs
}
