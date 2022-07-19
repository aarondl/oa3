{{- $.Import "net/http" -}}
{{- $.Import "net/http/httputil" -}}
{{- $.Import "context" -}}

type ctxKey string

const (
	ctxKeyDebug ctxKey = "debug"
)

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
		return Client{client: httpClient}
	}
	return Client{client: apiHTTPClient}
}

func (c Client) doRequest(req *http.Request) (*http.Response, error) {
	if c.limiter != nil {
		if err := c.limiter.Wait(req.Context(), 1); err != nil {
			return nil, err
		}
	}

	if hasDebug(req.Context()) {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, fmt.Errorf("failed to emit debugging info: %w", err)
		}
		fmt.Printf("%s\n", reqDump)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if hasDebug(req.Context()) {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to emit debugging info: %w", err)
		}
		fmt.Printf("%s\n", respDump)
	}

	return resp, nil
}
