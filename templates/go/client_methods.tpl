{{- $.Import "errors"}}
{{- range $url, $path := $.Spec.Paths -}}
    {{- range $method, $op := $path.Operations -}}
        {{- $opname := (camelcase $op.OperationID) -}}
        {{- $.Import "net/http"}}
// {{$opname}} {{$method}} {{$url}}
{{if $op.Description -}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" $op.Description)}}
{{end -}}
func (c {{$.Name}}) {{$opname}}(ctx context.Context
		{{- if $op.RequestBody -}}
			, {{$media := index $op.RequestBody.Content "application/json" -}}
			body{{" " -}}
				{{- if $media.Schema.Ref -}}
					{{- if not (isInlinePrimitive $media.Schema.Schema) -}}*{{- end -}}
					{{- refName $media.Schema.Ref -}}
				{{- else if isInlinePrimitive $media.Schema.Schema -}}
					{{- primitive $ $media.Schema.Schema $op.RequestBody.Required -}}
				{{- else -}}
					{{title $op.OperationID}}Inline
				{{- end -}}
		{{- end -}}
		{{- range $param := $op.Parameters -}}
		, {{untitle (camelcase $param.Name)}} {{primitive $ $param.Schema.Schema $param.Required -}}
		{{- end -}}
	) ({{title $op.OperationID}}Response, *http.Response, error) {
    urlStr := `{{$url}}`
	{{- range $param := $op.Parameters -}}
        {{- if and (eq $param.In "path") -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- $.Import "strings"}}
            urlStr = strings.Replace(urlStr, `{{"{"}}{{$param.Name}}{{"}"}}`, fmt.Sprintf("%v", {{$pname}}), 1) 
        {{- end -}}
	{{- end}}
	req := http.NewRequest(http.Method{{camelcase $method}}, urlStr, nil)
	{{- if $op.RequestBody}}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req.Body = bytes.NewReader(bodyBytes)
	{{- end -}}
	{{- range $param := $op.Parameters -}}
        {{- if and (eq $param.In "header") (not (has (lower $param.Name) (list "accept" "content-type" "authorization"))) -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- if or (not $param.Required) $param.Schema.Schema.Nullable }}
                if val, ok := {{$pname}}.Get(); ok {
                    req.Header().Set(`{{$param.Name}}`, fmt.Sprintf("%v", val))
                }
            {{- else}}
                req.Header().Set(`{{$param.Name}}`, fmt.Sprintf("%v", {{$pname}}))
            {{- end -}}
        {{- end -}}
	{{- end}}
    var queryStringValues url.Values
    {{$queryBuilt := false -}}
	{{- range $param := $op.Parameters -}}
        {{- if eq $param.In "query" -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- if not $queryBuilt -}}
                {{- $queryBuilt = true -}}
            if queryStringValues == nil {
                queryStringValues = make(url.Values)
            }
            {{- end -}}
            {{- if or (not $param.Required) $param.Schema.Schema.Nullable }}
                if val, ok := {{$pname}}.Get(); ok {
                    req.queryStringValues.Set(`{{$param.Name}}`, fmt.Sprintf("%v", val))
                }
            {{- else}}
                req.queryStringValues.Set(`{{$param.Name}}`, fmt.Sprintf("%v", {{$pname}}))
            {{- end -}}
        {{- end -}}
	{{- end}}
    if len(queryStringValues) > 0 {
        req.URL.RawQuery = queryStringValues.Encode()
    }

	{{range $param := $op.Parameters -}}
        {{- if eq $param.In "cookie" -}}
    // $param.Name cookie param not supported
        {{- end -}}
	{{- end -}}

    httpResp, err := c.doRequest(req)
    if err != nil {
        return nil, nil, err
    }

    var resp {{title $op.OperationID}}Response
    switch httpResp.Status {
    {{- $hasDefault := index $op.Responses "default"}}
    {{- range $code, $resp := $op.Responses}}
    {{if not (eq $code "default")}}case {{end -}}{{$code}}:
        {{$rkind := responseKind $op $code -}}
        {{- if eq $rkind "wrapped"}}
            var respObject {{$opname}}{{$code}}WrappedResponse
            {{- $.Import "io"}}
            {{- $.Import "json"}}
            b, err := io.ReadAll(resp.Body)
            if err != nil {
                return nil, nil, err
            }
            if err = json.Unmarshal(b, &respObject); err != nil {
                return nil, nil, err
            }
            resp = respObject
        {{- end -}}
    {{- end}}
    {{- if not $hasDefault}}
    default:
        return nil, nil, errors.Errorf("unknown response code")
    {{- end -}}
    }

    return resp, httpResp, nil
}
	{{end -}}
{{- end}}

{{$needHTTPStatuses := dict -}}

{{- range $url, $path := $.Spec.Paths -}}
    {{- range $method, $op := $path.Operations -}}
    {{- $opname := title $op.OperationID}}
// {{$opname}}Response one-of enforcer
//
// Implementors:
        {{- range $code, $resp := $op.Responses}}
// - {{if $resp.Headers -}}
        {{- if $resp.Content -}}
    {{$opname}}{{$code}}HeadersResponse
        {{- else -}}
            {{- $statusName := camelcase (httpStatus (atoi $code)) -}}
            {{- $_ := set $needHTTPStatuses $statusName "" -}}
    HTTPStatus{{$statusName}}
        {{- end -}}
     {{- else if $resp.Content -}}
            {{- $schema := index $resp.Content "application/json" -}}
            {{- $schema.Schema.Ref -}}
    {{- else -}}
                {{- $statusName := camelcase (httpStatus (atoi $code)) -}}
                {{- $_ := set $needHTTPStatuses $statusName "" -}}
    HTTPStatus{{$statusName}}
            {{- end -}}
        {{- end}}
type {{$opname}}Response interface {
    {{$opname}}Impl()
}

        {{- range $code, $resp := $op.Responses}}
            {{$rkind := responseKind $op $code -}}
            {{- if eq $rkind "wrapped" -}}
                {{- /* Headers, or duplicate responses produce a wrapper struct */}}
// {{$opname}}{{$code}}WrappedResponse wraps the normal body response with a
// struct to be able to additionally return headers or differentiate between
// multiple response codes with the same response body.
type {{$opname}}{{$code}}WrappedResponse struct {
                {{- range $hname, $header := $resp.Headers}}
    Header{{$hname | replace "-" "" | title}} {{if $header.Required -}}
                                    string
                                {{- else -}}
                                    {{- $.Import "github.com/aarondl/opt/omit" -}}
                                    omit.Val[string]
                                {{- end -}}
                {{- end -}}
    {{- $statusName := camelcase (httpStatus (atoi $code))}}
    Body {{if $resp.Content}}{{refName (index $resp.Content "application/json").Schema.Ref }}{{else}}HTTPStatus{{$statusName}}{{end}}
}

// {{$opname}}Impl implements {{$opname}}Response({{$code}}) for {{$opname}}{{$code}}WrappedResponse
func ({{$opname}}{{$code}}WrappedResponse) {{$opname}}Impl() {}

            {{- else if eq $rkind "empty" -}}
            {{- /* If there's no headers and no response body */ -}}
{{- $statusName := camelcase (httpStatus (atoi $code))}}
// {{$opname}}Impl implements {{$opname}}Response({{$code}}) for HTTPStatus{{$statusName}}
func (HTTPStatus{{$statusName}}) {{$opname}}Impl() {}
            {{- else -}}
                {{- /* If there's no headers */ -}}
                {{- $schema := index $resp.Content "application/json"}}
// {{$opname}}Impl implements {{$opname}}HeadersResponse({{$code}}) for {{refName $schema.Schema.Ref}}
func ({{if $schema.Schema.Ref}}{{refName $schema.Schema.Ref}}{{else}}{{title $opname}}{{title $code}}Inline{{end}}) {{$opname}}Impl() {}
            {{- end -}}
        {{- end -}}
    {{- end -}}
{{- end}}

{{range $status, $_ := $needHTTPStatuses -}}
// HTTPStatus{{$status}} is an empty response
type HTTPStatus{{$status}} struct {}
{{end -}}
