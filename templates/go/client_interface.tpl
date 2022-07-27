{{- $.Import "net/http" -}}
{{- $.Import "net/http/httputil" -}}
{{- $.Import "context" -}}
{{- $.Import "fmt" -}}
{{- $.Import "time" -}}

type ctxKey string

const (
	ctxKeyDebug ctxKey = "debug"
)

// BaseURLBuilder builds a base url. Implementations are likely simple fixed
// strings or slightly more complicated variable replacement strings with
// defaults.
//
// Implementors:
{{- range $server := $.Spec.Servers}}
// - {{filterNonIdentChars $server.URL | title}}
{{- end}}
type BaseURLBuilder interface {
	ToURL() string
}

// BaseURLSimple is a simple base url builder that's just a static string
type BaseURLSimple string
func (b BaseURLSimple) ToURL() string { return string(b) }

{{/* Called with an oa3 Server type */}}
{{define "server"}}
{{- $server := $.Object -}}
{{- $serverName := filterNonIdentChars $server.URL | title -}}
{{- if $server.Description}}
// {{wrapWith 70 "\n// " (trimSuffix "\n" $server.Description)}}
{{- end -}}
	{{- if $server.Variables}}
type {{$serverName}} struct {
		{{- range $k := (keysReflect $server.Variables | sortAlpha)}}
			{{$k | title}} string
		{{- end}}
}
		{{- with $.Name}}
func ({{$serverName}}) {{$.Name}}ToURL() {}
		{{- end}}
func (s {{$serverName}}) ToURL() string {
		uri := `{{$server.URL}}`
		{{- range $k := (keysReflect $server.Variables | sortAlpha) -}}
			{{- $varname := $k | title -}}
			{{- $v := index $server.Variables $k -}}
			{{- $.Import "strings"}}
		if len(s.{{$varname}}) != 0 {
			{{- with $enumVals := $v.Enum}}
			switch s.{{$varname}} {
			case {{range $i, $enum := $enumVals}}{{if ne $i 0}}, {{end}}`{{$enum}}`{{end}}:
			default:
				panic("unknown server variable enum value: " + s.{{$varname}})
			}
			{{- end}}
			uri = strings.ReplaceAll(uri, `{{"{"}}{{$k}}{{"}"}}`, s.{{$varname}})
		} else {
			uri = strings.ReplaceAll(uri, `{{"{"}}{{$k}}{{"}"}}`, `{{$v.Default}}`)
		}
		{{- end}}
		return uri
}
	{{- else}}
var {{$serverName}} = BaseURLSimple{{$.Name}}(`{{$server.URL}}`)
	{{- end -}}
{{- end -}}

{{- range $server := $.Spec.Servers -}}
{{template "server" (newData $ "" $server)}}
{{- end}}

{{- range $url, $path := $.Spec.Paths -}}
	{{- with $servers := $path.Servers -}}
		{{- $intfName := $url | filterNonIdentChars | title}}
// BaseURLSimple{{$intfName}} is a simple string url like BaseURLSimple
type BaseURLSimple{{$intfName}} string
func (b BaseURLSimple{{$intfName}}) ToURL() string { return string(b) }
func (b BaseURLSimple{{$intfName}}) {{$intfName}}Satisfy() {}

// BaseURLBuilder{{$intfName}} builds a base url like BaseURLBuilder but
// restricts the implementing types to a smaller subset.
//
// Implementors:
		{{- range $server := $servers}}
// - {{filterNonIdentChars $server.URL | title -}}
		{{- end}}
type BaseURLBuilder{{$intfName}} interface {
	BaseURLBuilder
	{{$intfName}}Satisfy()
}
		{{- range $server := $servers -}}
		{{template "server" (newData $ $intfName $server)}}
		{{- end}}
	{{- end -}}
    {{- range $method, $op := $path.Operations -}}
		{{- with $servers := $op.Servers -}}
		{{- $intfName := printf "%s%s" ($url | filterNonIdentChars | title) ($method | filterNonIdentChars | title) }}
// BaseURLSimple{{$intfName}} is a simple string url like BaseURLSimple
type BaseURLSimple{{$intfName}} string
func (b BaseURLSimple{{$intfName}}) ToURL() string { return string(b) }
func (b BaseURLSimple{{$intfName}}) {{$intfName}}Satisfy() {}

// BaseURLBuilder{{$intfName}} builds a base url like BaseURLBuilder but
// restricts the implementing types to a smaller subset.
//
// Implementors:
			{{- range $server := $servers}}
// - {{filterNonIdentChars $server.URL | title -}}
			{{- end}}
type BaseURLBuilder{{$intfName}} interface {
	BaseURLBuilder
	{{$intfName}}Satisfy()
}
			{{- with $servers := $op.Servers -}}
				{{- range $server := $servers -}}
			{{template "server" (newData $ $intfName $server)}}
				{{- end -}}
			{{- end -}}
		{{- end -}}
	{{- end -}}
{{- end}}

var (
	apiHTTPClient = &http.Client{Timeout: time.Second * 5}
)

// Client is a generated package for consuming an openapi spec.
{{- if $.Spec.Info.Description}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" $.Spec.Info.Description)}}
{{end -}}
type Client struct {
	httpClient *http.Client
	httpHandler http.Handler
	{{- $.Import "golang.org/x/time/rate"}}
	limiter *rate.Limiter
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
		{{- $.Import "net/http/httptest"}}
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

{{template "responses" $}}
