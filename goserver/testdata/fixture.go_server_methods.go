// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aarondl/chrono"
	"github.com/aarondl/oa3/support"
	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AlreadyHandled is an interface which an error return type can optionally
// implement to stop the generated method from responding in anyway, it will
// swallow the error and not touch the ResponseWriter if this method returns
// true.
type AlreadyHandled interface {
	AlreadyHandled() bool
}

// ErrHandled is a sentinel error that implements
// the AlreadyHandled interface which prevents the
// generated handler's response code from firing.
type ErrHandled struct{}

// Error implements error
func (ErrHandled) Error() string { return "already handled" }

// AlreadyHandled implements AlreadyHandled
func (ErrHandled) AlreadyHandled() bool { return true }

// authenticate post /auth
func (o GoServer) authenticateOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.Authenticate(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testarrayrequest get /test/array/request
func (o GoServer) testarrayrequestOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	var reqBody TestArrayRequestInline

	if r.Body == nil {
		return support.ErrNoBody
	} else {
		var buf *bytes.Buffer
		if buf, err = support.ReadJSONBuffer(r, &reqBody); err != nil {
			return err
		}

		defer support.ReturnJSONBuffer(buf)
		r.Body = io.NopCloser(buf)

		if newErrs := Validate(reqBody); newErrs != nil {
			errs = support.MergeErrs(errs, newErrs)
		}
	}
	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestArrayRequest(w, r, reqBody)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testmapsarrayinline get /test/arraymaps
func (o GoServer) testmapsarrayinlineOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestMapsArrayInline(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	if err := support.WriteJSON(w, ret); err != nil {
		return err
	}

	return nil
}

// testmapsarrayref post /test/arraymaps
func (o GoServer) testmapsarrayrefOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestMapsArrayRef(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	if err := support.WriteJSON(w, ret); err != nil {
		return err
	}

	return nil
}

// testenumqueryrequest get /test/enum/query/request
func (o GoServer) testenumqueryrequestOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs
	const n0 = `sort`
	query := r.URL.Query()
	s0 := query[n0]
	s0Exists := len(s0) > 0 && len(s0[0]) > 0
	var p0 TestEnumQueryRequestGetSortParam
	if !s0Exists || len(s0) == 0 {
		errs = support.AddErrs(errs, n0, errors.New(`must be provided and not be empty`))
	} else {
		p0 = TestEnumQueryRequestGetSortParam(s0[0])
		if newErrs := Validate(TestEnumQueryRequestGetSortParam(p0)); newErrs != nil {
			errs = support.AddErrsFlatten(errs, n0, newErrs)
		}

	}

	var reqBody TestEnumQueryRequestInline

	if r.Body == nil {
		return support.ErrNoBody
	} else {
		var buf *bytes.Buffer
		if buf, err = support.ReadJSONBuffer(r, &reqBody); err != nil {
			return err
		}

		defer support.ReturnJSONBuffer(buf)
		r.Body = io.NopCloser(buf)

		if newErrs := Validate(reqBody); newErrs != nil {
			errs = support.MergeErrs(errs, newErrs)
		}
	}
	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestEnumQueryRequest(w, r, reqBody, p0)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testinlineprimitivebody get /test/inline
func (o GoServer) testinlineprimitivebodyOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	var reqBody TestInlinePrimitiveBodyInline

	if r.Body == nil {
		return support.ErrNoBody
	} else {
		var buf *bytes.Buffer
		if buf, err = support.ReadJSONBuffer(r, &reqBody); err != nil {
			return err
		}

		defer support.ReturnJSONBuffer(buf)
		r.Body = io.NopCloser(buf)

		if newErrs := Validate(reqBody); newErrs != nil {
			errs = support.MergeErrs(errs, newErrs)
		}
	}
	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestInlinePrimitiveBody(w, r, string(reqBody))
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testinline post /test/inline
func (o GoServer) testinlineOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	var reqBody TestInlineInline

	if r.Body == nil {
		return support.ErrNoBody
	} else {
		var buf *bytes.Buffer
		if buf, err = support.ReadJSONBuffer(r, &reqBody); err != nil {
			return err
		}

		defer support.ReturnJSONBuffer(buf)
		r.Body = io.NopCloser(buf)

		if newErrs := Validate(reqBody); newErrs != nil {
			errs = support.MergeErrs(errs, newErrs)
		}
	}
	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestInline(w, r, reqBody)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	switch respBody := ret.(type) {

	case TestInline200Inline:
		w.WriteHeader(200)

		if err := support.WriteJSON(w, respBody); err != nil {
			return err
		}
	case *TestInline200Inline:
		w.WriteHeader(200)

		if err := support.WriteJSON(w, respBody); err != nil {
			return err
		}
	case TestInline201Inline:
		w.WriteHeader(201)

		if err := support.WriteJSON(w, respBody); err != nil {
			return err
		}
	case *TestInline201Inline:
		w.WriteHeader(201)

		if err := support.WriteJSON(w, respBody); err != nil {
			return err
		}
	default:
		_ = respBody
		panic("impossible case")
	}

	return nil
}

// testinlineresponsecomponent post /test/inlineresponsecomponent
func (o GoServer) testinlineresponsecomponentOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestInlineResponseComponent(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	if err := support.WriteJSON(w, ret); err != nil {
		return err
	}

	return nil
}

// testinlineresponsecomponentmultiple post /test/inlineresponsecomponentmultiple
func (o GoServer) testinlineresponsecomponentmultipleOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestInlineResponseComponentMultiple(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	switch respBody := ret.(type) {

	case InlineResponseTestInline:
		w.WriteHeader(200)

		if err := support.WriteJSON(w, respBody); err != nil {
			return err
		}
	case *InlineResponseTestInline:
		w.WriteHeader(200)

		if err := support.WriteJSON(w, respBody); err != nil {
			return err
		}
	case HTTPStatusCreated:
		w.WriteHeader(201)
	case *HTTPStatusCreated:
		w.WriteHeader(201)
	default:
		_ = respBody
		panic("impossible case")
	}

	return nil
}

// testmapsinline get /test/maps
func (o GoServer) testmapsinlineOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestMapsInline(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	if err := support.WriteJSON(w, ret); err != nil {
		return err
	}

	return nil
}

// testmapsref post /test/maps
func (o GoServer) testmapsrefOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestMapsRef(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	if err := support.WriteJSON(w, ret); err != nil {
		return err
	}

	return nil
}

// testqueryintarrayparam post /test/queryintarrayparam
func (o GoServer) testqueryintarrayparamOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs
	const n0 = `intarray`
	query := r.URL.Query()
	s0 := query[n0]
	s0Exists := len(s0) > 0 && len(s0[0]) > 0
	var p0 omit.Val[TestQueryIntArrayParamPostIntarrayParam]
	if s0Exists {
		c0, err := support.ExplodedFormArrayToSlice[int32](s0, support.StringToInt[int32])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n0, `TestQueryIntArrayParamPostIntarrayParam`, err)
		}
		p0.Set(c0)

	}
	const n1 = `intarrayrequired`
	s1 := query[n1]
	s1Exists := len(s1) > 0 && len(s1[0]) > 0
	var p1 TestQueryIntArrayParamPostIntarrayrequiredParam
	if !s1Exists || len(s1) == 0 {
		errs = support.AddErrs(errs, n1, errors.New(`must be provided and not be empty`))
	} else {
		c1, err := support.ExplodedFormArrayToSlice[int32](s1, support.StringToInt[int32])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n1, `TestQueryIntArrayParamPostIntarrayrequiredParam`, err)
		}
		p1 = c1

	}

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestQueryIntArrayParam(w, r, p0, p1)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testserverpathoverriderequest get /test/servers
func (o GoServer) testserverpathoverriderequestOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestServerPathOverrideRequest(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testserveropoverriderequest post /test/servers
func (o GoServer) testserveropoverriderequestOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestServerOpOverrideRequest(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testsingleserverpathoverriderequest get /test/single_servers
func (o GoServer) testsingleserverpathoverriderequestOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestSingleServerPathOverrideRequest(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testsingleserveropoverriderequest post /test/single_servers
func (o GoServer) testsingleserveropoverriderequestOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestSingleServerOpOverrideRequest(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testtypeoverrides get /test/type_overrides
func (o GoServer) testtypeoverridesOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs
	const n0 = `number`
	query := r.URL.Query()
	s0 := query[n0]
	s0Exists := len(s0) > 0 && len(s0[0]) > 0
	var p0 decimal.Decimal
	if !s0Exists || len(s0) == 0 {
		errs = support.AddErrs(errs, n0, errors.New(`must be provided and not be empty`))
	} else {
		c0, err := support.StringToDecimal(s0[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n0, `decimal.Decimal`, err)
		}
		p0 = c0

	}
	const n1 = `date`
	s1 := query[n1]
	s1Exists := len(s1) > 0 && len(s1[0]) > 0
	var p1 chrono.Date
	if !s1Exists || len(s1) == 0 {
		errs = support.AddErrs(errs, n1, errors.New(`must be provided and not be empty`))
	} else {
		c1, err := support.StringToChronoDate(s1[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n1, `chrono.Date`, err)
		}
		p1 = c1

	}
	const n2 = `number_null`
	s2 := query[n2]
	s2Exists := len(s2) > 0 && len(s2[0]) > 0
	var p2 null.Val[decimal.Decimal]
	if !s2Exists || len(s2) == 0 {
		errs = support.AddErrs(errs, n2, errors.New(`must be provided and not be empty`))
	} else {
		c2, err := support.StringToDecimal(s2[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n2, `decimal.Decimal`, err)
		}
		p2.Set(c2)

	}
	const n3 = `date_null`
	s3 := query[n3]
	s3Exists := len(s3) > 0 && len(s3[0]) > 0
	var p3 null.Val[chrono.Date]
	if !s3Exists || len(s3) == 0 {
		errs = support.AddErrs(errs, n3, errors.New(`must be provided and not be empty`))
	} else {
		c3, err := support.StringToChronoDate(s3[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n3, `chrono.Date`, err)
		}
		p3.Set(c3)

	}
	const n4 = `number_non_req`
	s4 := query[n4]
	s4Exists := len(s4) > 0 && len(s4[0]) > 0
	var p4 omit.Val[decimal.Decimal]
	if s4Exists {
		c4, err := support.StringToDecimal(s4[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n4, `decimal.Decimal`, err)
		}
		p4.Set(c4)

	}
	const n5 = `date_non_req`
	s5 := query[n5]
	s5Exists := len(s5) > 0 && len(s5[0]) > 0
	var p5 omit.Val[chrono.Date]
	if s5Exists {
		c5, err := support.StringToChronoDate(s5[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n5, `chrono.Date`, err)
		}
		p5.Set(c5)

	}

	var reqBody Primitives

	if r.Body == nil {
		return support.ErrNoBody
	} else {
		var buf *bytes.Buffer
		if buf, err = support.ReadJSONBuffer(r, &reqBody); err != nil {
			return err
		}

		defer support.ReturnJSONBuffer(buf)
		r.Body = io.NopCloser(buf)

		if newErrs := Validate(reqBody); newErrs != nil {
			errs = support.MergeErrs(errs, newErrs)
		}
	}
	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestTypeOverrides(w, r, &reqBody, p0, p1, p2, p3, p4, p5)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)

	return nil
}

// testunknownbodytype post /test/unknown/body/type
func (o GoServer) testunknownbodytypeOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.TestUnknownBodyType(w, r)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(200)
	if ret != nil {
		if _, err := io.Copy(w, ret); err != nil {
			return err
		}
		if err := ret.Close(); err != nil {
			return err
		}
	}

	return nil
}

// getuser get /users/{id}
func (o GoServer) getuserOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs
	const n0 = `id`
	s0, s0Exists := []string{chi.URLParam(r, n0)}, true
	var p0 string
	if !s0Exists || len(s0) == 0 {
		errs = support.AddErrs(errs, n0, errors.New(`must be provided and not be empty`))
	} else {
		p0 = s0[0]

	}
	const n1 = `param_component`
	query := r.URL.Query()
	s1 := query[n1]
	s1Exists := len(s1) > 0 && len(s1[0]) > 0
	var p1 string
	if !s1Exists || len(s1) == 0 {
		errs = support.AddErrs(errs, n1, errors.New(`must be provided and not be empty`))
	} else {
		p1 = s1[0]

	}
	const n2 = `valid_str`
	s2 := r.Header[http.CanonicalHeaderKey(n2)]
	s2Exists := len(s2) > 0 && len(s2[0]) > 0
	var p2 omitnull.Val[GetUserGetValidStrParam]
	if s2Exists {
		p2.Set(GetUserGetValidStrParam(s2[0]))
		if newErrs := Validate(GetUserGetValidStrParam(p2.GetOrZero())); newErrs != nil {
			errs = support.AddErrsFlatten(errs, n2, newErrs)
		}

	}
	const n3 = `req_valid_str`
	s3 := query[n3]
	s3Exists := len(s3) > 0 && len(s3[0]) > 0
	var p3 null.Val[GetUserGetReqValidStrParam]
	if !s3Exists || len(s3) == 0 {
		errs = support.AddErrs(errs, n3, errors.New(`must be provided and not be empty`))
	} else {
		p3.Set(GetUserGetReqValidStrParam(s3[0]))
		if newErrs := Validate(GetUserGetReqValidStrParam(p3.GetOrZero())); newErrs != nil {
			errs = support.AddErrsFlatten(errs, n3, newErrs)
		}

	}
	const n4 = `valid_int`
	s4 := query[n4]
	s4Exists := len(s4) > 0 && len(s4[0]) > 0
	var p4 omit.Val[int]
	if s4Exists {
		c4, err := support.StringToInt[int](s4[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n4, `int`, err)
		}
		p4.Set(c4)
		ers = nil
		if err := support.ValidateMaxNumber(p4.GetOrZero(), 5, true); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMinNumber(p4.GetOrZero(), 2, false); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMultipleOfInt(p4.GetOrZero(), 2); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			errs = support.AddErrs(errs, n4, ers...)
		}

	}
	const n5 = `req_valid_int`
	s5 := query[n5]
	s5Exists := len(s5) > 0 && len(s5[0]) > 0
	var p5 int
	if !s5Exists || len(s5) == 0 {
		errs = support.AddErrs(errs, n5, errors.New(`must be provided and not be empty`))
	} else {
		c5, err := support.StringToInt[int](s5[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n5, `int`, err)
		}
		p5 = c5
		ers = nil
		if err := support.ValidateMaxNumber(p5, 5, true); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMinNumber(p5, 2, false); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMultipleOfInt(p5, 2); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			errs = support.AddErrs(errs, n5, ers...)
		}

	}
	const n6 = `valid_num`
	s6 := query[n6]
	s6Exists := len(s6) > 0 && len(s6[0]) > 0
	var p6 omit.Val[float64]
	if s6Exists {
		c6, err := support.StringToFloat[float64](s6[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n6, `float64`, err)
		}
		p6.Set(c6)
		ers = nil
		if err := support.ValidateMaxNumber(p6.GetOrZero(), 10.5, false); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMinNumber(p6.GetOrZero(), 5.5, true); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMultipleOfFloat(p6.GetOrZero(), 2.5); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			errs = support.AddErrs(errs, n6, ers...)
		}

	}
	const n7 = `req_valid_num`
	s7 := query[n7]
	s7Exists := len(s7) > 0 && len(s7[0]) > 0
	var p7 float64
	if !s7Exists || len(s7) == 0 {
		errs = support.AddErrs(errs, n7, errors.New(`must be provided and not be empty`))
	} else {
		c7, err := support.StringToFloat[float64](s7[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n7, `float64`, err)
		}
		p7 = c7
		ers = nil
		if err := support.ValidateMaxNumber(p7, 10.5, false); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMinNumber(p7, 5.5, true); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMultipleOfFloat(p7, 2.5); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			errs = support.AddErrs(errs, n7, ers...)
		}

	}
	const n8 = `valid_bool`
	s8 := query[n8]
	s8Exists := len(s8) > 0 && len(s8[0]) > 0
	var p8 omit.Val[bool]
	if s8Exists {
		c8, err := support.StringToBool(s8[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n8, `bool`, err)
		}
		p8.Set(c8)

	}
	const n9 = `req_valid_bool`
	s9 := query[n9]
	s9Exists := len(s9) > 0 && len(s9[0]) > 0
	var p9 bool
	if !s9Exists || len(s9) == 0 {
		errs = support.AddErrs(errs, n9, errors.New(`must be provided and not be empty`))
	} else {
		c9, err := support.StringToBool(s9[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n9, `bool`, err)
		}
		p9 = c9

	}
	const n10 = `req_str_format`
	s10 := query[n10]
	s10Exists := len(s10) > 0 && len(s10[0]) > 0
	var p10 uuid.UUID
	if !s10Exists || len(s10) == 0 {
		errs = support.AddErrs(errs, n10, errors.New(`must be provided and not be empty`))
	} else {
		c10, err := support.StringToUUID(s10[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n10, `uuid.UUID`, err)
		}
		p10 = c10

	}
	const n11 = `date_time`
	s11 := query[n11]
	s11Exists := len(s11) > 0 && len(s11[0]) > 0
	var p11 chrono.DateTime
	if !s11Exists || len(s11) == 0 {
		errs = support.AddErrs(errs, n11, errors.New(`must be provided and not be empty`))
	} else {
		c11, err := support.StringToChronoDateTime(s11[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n11, `chrono.DateTime`, err)
		}
		p11 = c11

	}
	const n12 = `date`
	s12 := query[n12]
	s12Exists := len(s12) > 0 && len(s12[0]) > 0
	var p12 chrono.Date
	if !s12Exists || len(s12) == 0 {
		errs = support.AddErrs(errs, n12, errors.New(`must be provided and not be empty`))
	} else {
		c12, err := support.StringToChronoDate(s12[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n12, `chrono.Date`, err)
		}
		p12 = c12

	}
	const n13 = `time_val`
	s13 := query[n13]
	s13Exists := len(s13) > 0 && len(s13[0]) > 0
	var p13 chrono.Time
	if !s13Exists || len(s13) == 0 {
		errs = support.AddErrs(errs, n13, errors.New(`must be provided and not be empty`))
	} else {
		c13, err := support.StringToChronoTime(s13[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n13, `chrono.Time`, err)
		}
		p13 = c13

	}
	const n14 = `duration_val`
	s14 := query[n14]
	s14Exists := len(s14) > 0 && len(s14[0]) > 0
	var p14 time.Duration
	if !s14Exists || len(s14) == 0 {
		errs = support.AddErrs(errs, n14, errors.New(`must be provided and not be empty`))
	} else {
		c14, err := support.StringToDuration(s14[0])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n14, `time.Duration`, err)
		}
		p14 = c14

	}
	const n15 = `array_prim_explode`
	s15 := query[n15]
	s15Exists := len(s15) > 0 && len(s15[0]) > 0
	var p15 omit.Val[GetUserGetArrayPrimExplodeParam]
	if s15Exists {
		c15, err := support.ExplodedFormArrayToSlice[string](s15, support.StringNoOp)
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n15, `GetUserGetArrayPrimExplodeParam`, err)
		}
		p15.Set(c15)

	}
	const n16 = `array_prim_flat`
	s16 := query[n16]
	s16Exists := len(s16) > 0 && len(s16[0]) > 0
	var p16 GetUserGetArrayPrimFlatParam
	if !s16Exists || len(s16) == 0 {
		errs = support.AddErrs(errs, n16, errors.New(`must be provided and not be empty`))
	} else {
		c16, err := support.FlatFormArrayToSlice[string](s16, support.StringNoOp)
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n16, `GetUserGetArrayPrimFlatParam`, err)
		}
		p16 = c16

	}
	const n17 = `array_prim_int_explode`
	s17 := query[n17]
	s17Exists := len(s17) > 0 && len(s17[0]) > 0
	var p17 omit.Val[GetUserGetArrayPrimIntExplodeParam]
	if s17Exists {
		c17, err := support.ExplodedFormArrayToSlice[int](s17, support.StringToInt[int])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n17, `GetUserGetArrayPrimIntExplodeParam`, err)
		}
		p17.Set(c17)

	}
	const n18 = `array_prim_int_flat`
	s18 := query[n18]
	s18Exists := len(s18) > 0 && len(s18[0]) > 0
	var p18 GetUserGetArrayPrimIntFlatParam
	if !s18Exists || len(s18) == 0 {
		errs = support.AddErrs(errs, n18, errors.New(`must be provided and not be empty`))
	} else {
		c18, err := support.FlatFormArrayToSlice[int](s18, support.StringToInt[int])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n18, `GetUserGetArrayPrimIntFlatParam`, err)
		}
		p18 = c18

	}
	const n19 = `array_enum_explode`
	s19 := query[n19]
	s19Exists := len(s19) > 0 && len(s19[0]) > 0
	var p19 omit.Val[GetUserGetArrayEnumExplodeParam]
	if s19Exists {
		c19, err := support.ExplodedFormArrayToSlice[GetUserGetArrayEnumExplodeParamItem](s19, support.StringToString[string, GetUserGetArrayEnumExplodeParamItem])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n19, `GetUserGetArrayEnumExplodeParam`, err)
		}
		p19.Set(c19)

	}
	const n20 = `array_enum_flat`
	s20 := query[n20]
	s20Exists := len(s20) > 0 && len(s20[0]) > 0
	var p20 GetUserGetArrayEnumFlatParam
	if !s20Exists || len(s20) == 0 {
		errs = support.AddErrs(errs, n20, errors.New(`must be provided and not be empty`))
	} else {
		c20, err := support.FlatFormArrayToSlice[GetUserGetArrayEnumFlatParamItem](s20, support.StringToString[string, GetUserGetArrayEnumFlatParamItem])
		if err != nil {
			return fmt.Errorf("failed to convert parameter %q to %q: %w", n20, `GetUserGetArrayEnumFlatParam`, err)
		}
		p20 = c20

	}

	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.GetUser(w, r, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15, p16, p17, p18, p19, p20)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	_ = ret
	w.WriteHeader(304)

	return nil
}

// setuser post /users/{id}
func (o GoServer) setuserOp(w http.ResponseWriter, r *http.Request) error {
	var err error
	var ers []error
	var errs map[string][]string
	_, _, _ = err, ers, errs
	const n0 = `id`
	s0, s0Exists := []string{chi.URLParam(r, n0)}, true
	var p0 string
	if !s0Exists || len(s0) == 0 {
		errs = support.AddErrs(errs, n0, errors.New(`must be provided and not be empty`))
	} else {
		p0 = s0[0]

	}
	const n1 = `param_component`
	query := r.URL.Query()
	s1 := query[n1]
	s1Exists := len(s1) > 0 && len(s1[0]) > 0
	var p1 string
	if !s1Exists || len(s1) == 0 {
		errs = support.AddErrs(errs, n1, errors.New(`must be provided and not be empty`))
	} else {
		p1 = s1[0]

	}

	var reqBody Primitives

	if r.Body == nil {
		return support.ErrNoBody
	} else {
		var buf *bytes.Buffer
		if buf, err = support.ReadJSONBuffer(r, &reqBody); err != nil {
			return err
		}

		defer support.ReturnJSONBuffer(buf)
		r.Body = io.NopCloser(buf)

		if newErrs := Validate(reqBody); newErrs != nil {
			errs = support.MergeErrs(errs, newErrs)
		}
	}
	if errs != nil {
		return o.converter(errs)
	}

	ret, err := o.impl.SetUser(w, r, &reqBody, p0, p1)
	if err != nil {
		if alreadyHandled, ok := err.(AlreadyHandled); ok {
			if alreadyHandled.AlreadyHandled() {
				return nil
			}
		}
		return err
	}

	switch respBody := ret.(type) {

	case SetUserWrappedResponse:
		headers := w.Header()
		if val, ok := respBody.HeaderXResponseHeader.Get(); ok {
			headers.Set("X-Response-Header", val)
		}
		w.WriteHeader(respBody.Code)
		if err := support.WriteJSON(w, respBody.Body); err != nil {
			return err
		}
	case *SetUserWrappedResponse:
		headers := w.Header()
		if val, ok := respBody.HeaderXResponseHeader.Get(); ok {
			headers.Set("X-Response-Header", val)
		}
		w.WriteHeader(respBody.Code)
		if err := support.WriteJSON(w, respBody.Body); err != nil {
			return err
		}
	default:
		_ = respBody
		panic("impossible case")
	}

	return nil
}
