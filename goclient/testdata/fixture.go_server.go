// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/aarondl/opt/omit"
	"golang.org/x/time/rate"
)

type ctxKey string

const (
	ctxKeyDebug ctxKey = "debug"
)

// BaseURLBuilder builds a base url. Implementations are likely simple fixed
// strings or slightly more complicated variable replacement strings with
// defaults.
//
// Implementors:
// - Httpdevlocal
// - Httpprodlocalonetwo
// - Httpvariableslocalvariable
type BaseURLBuilder interface {
	ToURL() string
}

// BaseURLSimple is a simple base url builder that's just a static string
type BaseURLSimple string

func (b BaseURLSimple) ToURL() string { return string(b) }

// Local development
var Httpdevlocal = BaseURLSimple(`http://dev.local:3030`)

// Production
type Httpprodlocalonetwo struct {
	One string
	Two string
}

func (s Httpprodlocalonetwo) ToURL() string {
	uri := `http://prod.local:3030/{one}/{two}`
	if len(s.One) != 0 {
		uri = strings.ReplaceAll(uri, `{one}`, s.One)
	} else {
		uri = strings.ReplaceAll(uri, `{one}`, `one`)
	}
	if len(s.Two) != 0 {
		uri = strings.ReplaceAll(uri, `{two}`, s.Two)
	} else {
		uri = strings.ReplaceAll(uri, `{two}`, `two`)
	}
	return uri
}

// Variable path
type Httpvariableslocalvariable struct {
	Variable string
}

func (s Httpvariableslocalvariable) ToURL() string {
	uri := `http://variables.local:3030/{variable}`
	if len(s.Variable) != 0 {
		switch s.Variable {
		case `v1`, `v2`, `v3`:
		default:
			panic("unknown server variable enum value: " + s.Variable)
		}
		uri = strings.ReplaceAll(uri, `{variable}`, s.Variable)
	} else {
		uri = strings.ReplaceAll(uri, `{variable}`, `v1`)
	}
	return uri
}

// BaseURLSimpleTestservers is a simple string url like BaseURLSimple
type BaseURLSimpleTestservers string

func (b BaseURLSimpleTestservers) ToURL() string       { return string(b) }
func (b BaseURLSimpleTestservers) TestserversSatisfy() {}

// BaseURLBuilderTestservers builds a base url like BaseURLBuilder but
// restricts the implementing types to a smaller subset.
//
// Implementors:
// - Httppathdevlocal
// - Httppathprodlocalonetwo
// - Httppathvariableslocalvariable
type BaseURLBuilderTestservers interface {
	BaseURLBuilder
	TestserversSatisfy()
}

// Local development
var Httppathdevlocal = BaseURLSimpleTestservers(`http://path.dev.local:3030`)

// Production
type Httppathprodlocalonetwo struct {
	One string
	Two string
}

func (Httppathprodlocalonetwo) TestserversToURL() {}
func (s Httppathprodlocalonetwo) ToURL() string {
	uri := `http://path.prod.local:3030/{one}/{two}`
	if len(s.One) != 0 {
		uri = strings.ReplaceAll(uri, `{one}`, s.One)
	} else {
		uri = strings.ReplaceAll(uri, `{one}`, `one`)
	}
	if len(s.Two) != 0 {
		uri = strings.ReplaceAll(uri, `{two}`, s.Two)
	} else {
		uri = strings.ReplaceAll(uri, `{two}`, `two`)
	}
	return uri
}

// Variable path
type Httppathvariableslocalvariable struct {
	Variable string
}

func (Httppathvariableslocalvariable) TestserversToURL() {}
func (s Httppathvariableslocalvariable) ToURL() string {
	uri := `http://path.variables.local:3030/{variable}`
	if len(s.Variable) != 0 {
		switch s.Variable {
		case `v1`, `v2`, `v3`:
		default:
			panic("unknown server variable enum value: " + s.Variable)
		}
		uri = strings.ReplaceAll(uri, `{variable}`, s.Variable)
	} else {
		uri = strings.ReplaceAll(uri, `{variable}`, `v1`)
	}
	return uri
}

// BaseURLSimpleTestserversPost is a simple string url like BaseURLSimple
type BaseURLSimpleTestserversPost string

func (b BaseURLSimpleTestserversPost) ToURL() string           { return string(b) }
func (b BaseURLSimpleTestserversPost) TestserversPostSatisfy() {}

// BaseURLBuilderTestserversPost builds a base url like BaseURLBuilder but
// restricts the implementing types to a smaller subset.
//
// Implementors:
// - Httpopdevlocal
// - Httpopprodlocalonetwo
// - Httpopvariableslocalvariable
type BaseURLBuilderTestserversPost interface {
	BaseURLBuilder
	TestserversPostSatisfy()
}

// Local development
var Httpopdevlocal = BaseURLSimpleTestserversPost(`http://op.dev.local:3030`)

// Production
type Httpopprodlocalonetwo struct {
	One string
	Two string
}

func (Httpopprodlocalonetwo) TestserversPostToURL() {}
func (s Httpopprodlocalonetwo) ToURL() string {
	uri := `http://op.prod.local:3030/{one}/{two}`
	if len(s.One) != 0 {
		uri = strings.ReplaceAll(uri, `{one}`, s.One)
	} else {
		uri = strings.ReplaceAll(uri, `{one}`, `one`)
	}
	if len(s.Two) != 0 {
		uri = strings.ReplaceAll(uri, `{two}`, s.Two)
	} else {
		uri = strings.ReplaceAll(uri, `{two}`, `two`)
	}
	return uri
}

// Variable path
type Httpopvariableslocalvariable struct {
	Variable string
}

func (Httpopvariableslocalvariable) TestserversPostToURL() {}
func (s Httpopvariableslocalvariable) ToURL() string {
	uri := `http://op.variables.local:3030/{variable}`
	if len(s.Variable) != 0 {
		switch s.Variable {
		case `v1`, `v2`, `v3`:
		default:
			panic("unknown server variable enum value: " + s.Variable)
		}
		uri = strings.ReplaceAll(uri, `{variable}`, s.Variable)
	} else {
		uri = strings.ReplaceAll(uri, `{variable}`, `v1`)
	}
	return uri
}

var (
	apiHTTPClient = &http.Client{Timeout: time.Second * 5}
)

// Client is a generated package for consuming an openapi spec.
//
// A great api
type Client struct {
	httpClient  *http.Client
	httpHandler http.Handler
	limiter     *rate.Limiter
}

// WithDebug creates a context that will emit debugging information to stdout
// for each request.
func WithDebug(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyDebug, "t")
}

func hasDebug(ctx context.Context) bool {
	v := ctx.Value(ctxKeyDebug)
	return v != nil && v.(string) == "t"
}

// NewClient constructs an api client, optionally using a supplied http.Client
// to be able to add instrumentation or customized timeouts.
//
// If nil is supplied then this package's generated apiHTTPClient is used which
// has reasonable defaults for timeouts.
//
// It also takes an optional rate limiter to implement rate limiting.
func NewClient(httpClient *http.Client, limiter *rate.Limiter) Client {
	if httpClient != nil {
		return Client{httpClient: httpClient}
	}
	return Client{httpClient: apiHTTPClient}
}

// NewLocalClient constructs an api client, but takes in a handler to call
// with the prepared requests instead of an http client that will touch the
// network. Useful for testing.
func NewLocalClient(httpHandler http.Handler) Client {
	return Client{httpHandler: httpHandler}
}

func (c Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	if hasDebug(ctx) {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, fmt.Errorf("failed to emit debugging info: %w", err)
		}
		fmt.Printf("%s\n", reqDump)
	}

	var resp *http.Response
	if c.httpHandler != nil {
		w := httptest.NewRecorder()
		c.httpHandler.ServeHTTP(w, req)
		resp = w.Result()
	} else {
		var err error
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
	}

	if hasDebug(ctx) {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to emit debugging info: %w", err)
		}
		fmt.Printf("%s\n", respDump)
	}

	return resp, nil
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
