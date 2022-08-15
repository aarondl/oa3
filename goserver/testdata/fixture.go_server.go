// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
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

// Interface is the interface that an application server must implement
// in order to use this package.
//
// A great api
type Interface interface {
	// Authenticate post /auth
	Authenticate(w http.ResponseWriter, r *http.Request) (AuthenticateResponse, error)
	// TestArrayRequest get /test/array/request
	TestArrayRequest(w http.ResponseWriter, r *http.Request, body TestArrayRequestInline) (TestArrayRequestResponse, error)
	// TestEnumQueryRequest get /test/enum/query/request
	TestEnumQueryRequest(w http.ResponseWriter, r *http.Request, body TestEnumQueryRequestInline, sort TestEnumQueryRequestGetSortParam) (TestEnumQueryRequestResponse, error)
	// TestInlinePrimitiveBody get /test/inline
	TestInlinePrimitiveBody(w http.ResponseWriter, r *http.Request, body string) (TestInlinePrimitiveBodyResponse, error)
	// TestInline post /test/inline
	TestInline(w http.ResponseWriter, r *http.Request, body TestInlineInline) (TestInlineResponse, error)
	// TestServerPathOverrideRequest get /test/servers
	TestServerPathOverrideRequest(w http.ResponseWriter, r *http.Request) (TestServerPathOverrideRequestResponse, error)
	// TestServerOpOverrideRequest post /test/servers
	TestServerOpOverrideRequest(w http.ResponseWriter, r *http.Request) (TestServerOpOverrideRequestResponse, error)
	// TestTypeOverrides get /test/type_overrides
	TestTypeOverrides(w http.ResponseWriter, r *http.Request, body *Primitives, number decimal.Decimal, date chrono.Date, numberNull null.Val[decimal.Decimal], dateNull null.Val[chrono.Date], numberNonReq omit.Val[decimal.Decimal], dateNonReq omit.Val[chrono.Date]) (TestTypeOverridesResponse, error)
	// TestUnknownBodyType post /test/unknown/body/type
	TestUnknownBodyType(w http.ResponseWriter, r *http.Request) (TestUnknownBodyTypeResponse, error)
	// GetUser get /users/{id}
	// Retrieves a user with a long description that spans multiple lines so
	// that we can see that both wrapping and long-line support is not
	// bleeding over the sacred 80 char limit.
	GetUser(w http.ResponseWriter, r *http.Request, id string, validStr omitnull.Val[GetUserGetValidStrParam], reqValidStr null.Val[GetUserGetReqValidStrParam], validInt omit.Val[int], reqValidInt int, validNum omit.Val[float64], reqValidNum float64, validBool omit.Val[bool], reqValidBool bool, reqStrFormat uuid.UUID, dateTime chrono.DateTime, date chrono.Date, timeVal chrono.Time, durationVal time.Duration, arrayPrimExplode omit.Val[GetUserGetArrayPrimExplodeParam], arrayPrimFlat GetUserGetArrayPrimFlatParam, arrayPrimIntExplode omit.Val[GetUserGetArrayPrimIntExplodeParam], arrayPrimIntFlat GetUserGetArrayPrimIntFlatParam, arrayEnumExplode omit.Val[GetUserGetArrayEnumExplodeParam], arrayEnumFlat GetUserGetArrayEnumFlatParam) (GetUserResponse, error)
	// SetUser post /users/{id}
	// Sets a user
	SetUser(w http.ResponseWriter, r *http.Request, body *Primitives) (SetUserResponse, error)
}

type (
	// GoServer implements all the wrapper functionality for the API
	GoServer struct {
		impl      Interface
		converter support.ValidationConverter
		router    *chi.Mux
	}
)

// NewGoServer constructor
func NewGoServer(
	impl Interface,
	cnv support.ValidationConverter,
	eh support.ErrorHandler,
	mw support.MW,
) http.Handler {

	o := GoServer{
		impl:      impl,
		converter: cnv,
		router:    chi.NewRouter(),
	}

	// Untagged operations
	o.router.Group(func(r chi.Router) {
		if m, ok := mw[""]; ok {
			r.Use(m...)
		}
		r.Method(http.MethodPost, `/auth`, eh.Wrap(o.authenticateOp))
		r.Method(http.MethodGet, `/test/array/request`, eh.Wrap(o.testarrayrequestOp))
		r.Method(http.MethodGet, `/test/enum/query/request`, eh.Wrap(o.testenumqueryrequestOp))
		r.Method(http.MethodPost, `/test/inline`, eh.Wrap(o.testinlineOp))
		r.Method(http.MethodGet, `/test/inline`, eh.Wrap(o.testinlineprimitivebodyOp))
		r.Method(http.MethodPost, `/test/servers`, eh.Wrap(o.testserveropoverriderequestOp))
		r.Method(http.MethodGet, `/test/servers`, eh.Wrap(o.testserverpathoverriderequestOp))
		r.Method(http.MethodGet, `/test/type_overrides`, eh.Wrap(o.testtypeoverridesOp))
		r.Method(http.MethodPost, `/test/unknown/body/type`, eh.Wrap(o.testunknownbodytypeOp))
	})
	// users tagged operations
	o.router.Group(func(r chi.Router) {
		if m, ok := mw["users"]; ok {
			r.Use(m...)
		}
		r.Method(http.MethodGet, `/users/{id}`, eh.Wrap(o.getuserOp))
		r.Method(http.MethodPost, `/users/{id}`, eh.Wrap(o.setuserOp))
	})

	return o
}

// ServeHTTP implements http.Handler
func (o GoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.router.ServeHTTP(w, r)
}

type validatable interface {
	validateSchema() support.Errors
}

// Validate a schema
func Validate(toValidate validatable) support.Errors {
	return toValidate.validateSchema()
}

// AuthenticateResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type AuthenticateResponse interface {
	AuthenticateImpl()
}

// AuthenticateImpl implements AuthenticateResponse(200) for HTTPStatusOk
func (HTTPStatusOk) AuthenticateImpl() {}

// TestArrayRequestResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type TestArrayRequestResponse interface {
	TestArrayRequestImpl()
}

// TestArrayRequestImpl implements TestArrayRequestResponse(200) for HTTPStatusOk
func (HTTPStatusOk) TestArrayRequestImpl() {}

// TestEnumQueryRequestResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type TestEnumQueryRequestResponse interface {
	TestEnumQueryRequestImpl()
}

// TestEnumQueryRequestImpl implements TestEnumQueryRequestResponse(200) for HTTPStatusOk
func (HTTPStatusOk) TestEnumQueryRequestImpl() {}

// TestInlinePrimitiveBodyResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type TestInlinePrimitiveBodyResponse interface {
	TestInlinePrimitiveBodyImpl()
}

// TestInlinePrimitiveBodyImpl implements TestInlinePrimitiveBodyResponse(200) for HTTPStatusOk
func (HTTPStatusOk) TestInlinePrimitiveBodyImpl() {}

// TestInlineResponse one-of enforcer
//
// Implementors:
// - TestInline200Inline
// - TestInline201Inline
type TestInlineResponse interface {
	TestInlineImpl()
}

// TestInlineImpl implements TestInlineHeadersResponse(200) for
func (TestInline200Inline) TestInlineImpl() {}

// TestInlineImpl implements TestInlineHeadersResponse(201) for
func (TestInline201Inline) TestInlineImpl() {}

// TestServerPathOverrideRequestResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type TestServerPathOverrideRequestResponse interface {
	TestServerPathOverrideRequestImpl()
}

// TestServerPathOverrideRequestImpl implements TestServerPathOverrideRequestResponse(200) for HTTPStatusOk
func (HTTPStatusOk) TestServerPathOverrideRequestImpl() {}

// TestServerOpOverrideRequestResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type TestServerOpOverrideRequestResponse interface {
	TestServerOpOverrideRequestImpl()
}

// TestServerOpOverrideRequestImpl implements TestServerOpOverrideRequestResponse(200) for HTTPStatusOk
func (HTTPStatusOk) TestServerOpOverrideRequestImpl() {}

// TestTypeOverridesResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type TestTypeOverridesResponse interface {
	TestTypeOverridesImpl()
}

// TestTypeOverridesImpl implements TestTypeOverridesResponse(200) for HTTPStatusOk
func (HTTPStatusOk) TestTypeOverridesImpl() {}

// TestUnknownBodyTypeResponse one-of enforcer
//
// Implementors:
// - HTTPStatusOk
type TestUnknownBodyTypeResponse interface {
	TestUnknownBodyTypeImpl()
}

// TestUnknownBodyTypeImpl implements TestUnknownBodyTypeResponse(200) for HTTPStatusOk
func (HTTPStatusOk) TestUnknownBodyTypeImpl() {}

// GetUserResponse one-of enforcer
//
// Implementors:
// - HTTPStatusNotModified
type GetUserResponse interface {
	GetUserImpl()
}

// GetUserImpl implements GetUserResponse(304) for HTTPStatusNotModified
func (HTTPStatusNotModified) GetUserImpl() {}

// SetUserResponse one-of enforcer
//
// Implementors:
// - SetUser200HeadersResponse
// - Primitives - #/components/schemas/Primitives
type SetUserResponse interface {
	SetUserImpl()
}

// SetUser200WrappedResponse wraps the normal body response with a
// struct to be able to additionally return headers or differentiate between
// multiple response codes with the same response body.
type SetUser200WrappedResponse struct {
	HeaderXResponseHeader omit.Val[string]
	Body                  Primitives
}

// SetUserImpl implements SetUserResponse(200) for SetUser200WrappedResponse
func (SetUser200WrappedResponse) SetUserImpl() {}

// SetUserdefaultWrappedResponse wraps the normal body response with a
// struct to be able to additionally return headers or differentiate between
// multiple response codes with the same response body.
type SetUserdefaultWrappedResponse struct {
	Body Primitives
}

// SetUserImpl implements SetUserResponse(default) for SetUserdefaultWrappedResponse
func (SetUserdefaultWrappedResponse) SetUserImpl() {}

// HTTPStatusNotModified is an empty response
type HTTPStatusNotModified struct{}

// HTTPStatusOk is an empty response
type HTTPStatusOk struct{}

/*
Here is a copy pastable list of function signatures
for implementing the main interface

// Authenticate post /auth
func (a API) Authenticate(w http.ResponseWriter, r *http.Request) (oa3gen.AuthenticateResponse, error) {
    panic("not implemented")
}
// TestArrayRequest get /test/array/request
func (a API) TestArrayRequest(w http.ResponseWriter, r *http.Request, body oa3gen.TestArrayRequestInline) (oa3gen.TestArrayRequestResponse, error) {
    panic("not implemented")
}
// TestEnumQueryRequest get /test/enum/query/request
func (a API) TestEnumQueryRequest(w http.ResponseWriter, r *http.Request, body oa3gen.TestEnumQueryRequestInline, sort TestEnumQueryRequestGetSortParam) (oa3gen.TestEnumQueryRequestResponse, error) {
    panic("not implemented")
}
// TestInlinePrimitiveBody get /test/inline
func (a API) TestInlinePrimitiveBody(w http.ResponseWriter, r *http.Request, body string) (oa3gen.TestInlinePrimitiveBodyResponse, error) {
    panic("not implemented")
}
// TestInline post /test/inline
func (a API) TestInline(w http.ResponseWriter, r *http.Request, body oa3gen.TestInlineInline) (oa3gen.TestInlineResponse, error) {
    panic("not implemented")
}
// TestServerPathOverrideRequest get /test/servers
func (a API) TestServerPathOverrideRequest(w http.ResponseWriter, r *http.Request) (oa3gen.TestServerPathOverrideRequestResponse, error) {
    panic("not implemented")
}
// TestServerOpOverrideRequest post /test/servers
func (a API) TestServerOpOverrideRequest(w http.ResponseWriter, r *http.Request) (oa3gen.TestServerOpOverrideRequestResponse, error) {
    panic("not implemented")
}
// TestTypeOverrides get /test/type_overrides
func (a API) TestTypeOverrides(w http.ResponseWriter, r *http.Request, body *oa3gen.Primitives, number decimal.Decimal, date chrono.Date, numberNull null.Val[decimal.Decimal], dateNull null.Val[chrono.Date], numberNonReq omit.Val[decimal.Decimal], dateNonReq omit.Val[chrono.Date]) (oa3gen.TestTypeOverridesResponse, error) {
    panic("not implemented")
}
// TestUnknownBodyType post /test/unknown/body/type
func (a API) TestUnknownBodyType(w http.ResponseWriter, r *http.Request) (oa3gen.TestUnknownBodyTypeResponse, error) {
    panic("not implemented")
}
// GetUser get /users/{id}
// Retrieves a user with a long description that spans multiple lines so
// that we can see that both wrapping and long-line support is not
// bleeding over the sacred 80 char limit.
func (a API) GetUser(w http.ResponseWriter, r *http.Request, id string, validStr omitnull.Val[GetUserGetValidStrParam], reqValidStr null.Val[GetUserGetReqValidStrParam], validInt omit.Val[int], reqValidInt int, validNum omit.Val[float64], reqValidNum float64, validBool omit.Val[bool], reqValidBool bool, reqStrFormat uuid.UUID, dateTime chrono.DateTime, date chrono.Date, timeVal chrono.Time, durationVal time.Duration, arrayPrimExplode omit.Val[GetUserGetArrayPrimExplodeParam], arrayPrimFlat GetUserGetArrayPrimFlatParam, arrayPrimIntExplode omit.Val[GetUserGetArrayPrimIntExplodeParam], arrayPrimIntFlat GetUserGetArrayPrimIntFlatParam, arrayEnumExplode omit.Val[GetUserGetArrayEnumExplodeParam], arrayEnumFlat GetUserGetArrayEnumFlatParam) (oa3gen.GetUserResponse, error) {
    panic("not implemented")
}
// SetUser post /users/{id}
// Sets a user
func (a API) SetUser(w http.ResponseWriter, r *http.Request, body *oa3gen.Primitives) (oa3gen.SetUserResponse, error) {
    panic("not implemented")
}
*/
