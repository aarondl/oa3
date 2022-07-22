{{- range $url, $path := $.Spec.Paths -}}
    {{- range $method, $op := $path.Operations -}}
        {{- $opname := (title $op.OperationID) -}}
        {{- $.Import "net/http"}}
// {{$opname}} {{$method}} {{$url}}
{{- if $op.Description}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" $op.Description)}}
{{- end -}}
{{- $.Import "context"}}
func (_c Client) {{$opname}}(ctx context.Context
		{{- if $op.RequestBody -}}
			, {{$media := index $op.RequestBody.Content "application/json" -}}
			body{{" " -}}
				{{- if $media.Schema.Ref -}}
					{{- if not (isInlinePrimitive $media.Schema.Schema) -}}*{{- end -}}
					{{- refName $media.Schema.Ref -}}
                {{- else if not (or (eq $media.Schema.Schema.Type "object") (eq $media.Schema.Schema.Type "array")) -}}
					{{- primitive $ $media.Schema.Schema $op.RequestBody.Required -}}
				{{- else -}}
					{{title $op.OperationID}}Inline
				{{- end -}}
		{{- end -}}
		{{- range $param := $op.Parameters -}}
		, {{untitle (camelcase $param.Name)}}{{" "}}
            {{- if and $param.Schema.Schema.Enum (gt (len $param.Schema.Schema.Enum) 0) -}}
                {{omitnullWrap $ $param.Schema.Schema (printf "%s%sParam" ($op.OperationID | snakeToCamel | title) ($param.Name | snakeToCamel | title)) $param.Schema.Schema.Nullable $param.Required}}
            {{- else -}}
                {{- primitive $ $param.Schema.Schema $param.Required -}}
            {{- end -}}
		{{- end -}}
	) ({{title $op.OperationID}}Response, *http.Response, error) {
    _urlStr := `{{$url}}`
	{{- range $param := $op.Parameters -}}
        {{- if and (eq $param.In "path") -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- $.Import "strings"}}
            _urlStr = strings.Replace(_urlStr, `{{"{"}}{{$param.Name}}{{"}"}}`, fmt.Sprintf("%v", {{$pname}}), 1) 
        {{- end -}}
	{{- end}}
	_req, _err := http.NewRequestWithContext(ctx, http.Method{{camelcase $method}}, _urlStr, nil)
    if _err != nil {
        return nil, nil, _err
    }
	{{- if $op.RequestBody -}}
    {{- $.Import "encoding/json" }}
    {{- $.Import "bytes" }}
    {{- $.Import "io" }}
	_bodyBytes, _err := json.Marshal(body)
	if _err != nil {
		return nil, nil, _err
	}
	_req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
	{{- end -}}
	{{- range $param := $op.Parameters -}}
        {{- if and (eq $param.In "header") (not (has (lower $param.Name) (list "accept" "content-type" "authorization"))) -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- if or (not $param.Required) $param.Schema.Schema.Nullable }}
                if val, ok := {{$pname}}.Get(); ok {
                    _req.Header.Set(`{{$param.Name}}`, fmt.Sprintf("%v", val))
                }
            {{- else}}
                _req.Header.Set(`{{$param.Name}}`, fmt.Sprintf("%v", {{$pname}}))
            {{- end -}}
        {{- end -}}
	{{- end -}}
    {{- $.Import "net/url" }}
    var _query url.Values
    {{- $queryBuilt := false -}}
	{{- range $param := $op.Parameters -}}
        {{- if eq $param.In "query" -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- if not $queryBuilt -}}
                {{- $queryBuilt = true}}
            if _query == nil {
                _query = make(url.Values)
            }
            {{- end -}}
            {{- if or (not $param.Required) $param.Schema.Schema.Nullable }}
                if val, ok := {{$pname}}.Get(); ok {
                    _query.Set(`{{$param.Name}}`, fmt.Sprintf("%v", val))
                }
            {{- else}}
                _query.Set(`{{$param.Name}}`, fmt.Sprintf("%v", {{$pname}}))
            {{- end -}}
        {{- end -}}
	{{- end}}
    if len(_query) > 0 {
        _req.URL.RawQuery = _query.Encode()
    }

	{{range $param := $op.Parameters -}}
        {{- if eq $param.In "cookie" -}}
    // $param.Name cookie param not supported
        {{- end -}}
	{{- end -}}

    _httpResp, _err := _c.doRequest(ctx, _req)
    if _err != nil {
        return nil, nil, _err
    }

    var _resp {{title $op.OperationID}}Response
    switch _httpResp.Status {
    {{- $hasDefault := index $op.Responses "default"}}
    {{- range $code, $resp := $op.Responses}}
    {{if not (eq $code "default")}}case `{{$code}}`:{{- else -}}default:{{- end}}
        {{- $rkind := responseKind $op $code -}}
        {{- if eq $rkind "wrapped"}}
            var _respObject {{$opname}}{{$code}}WrappedResponse
            {{- $.Import "io"}}
            {{- $.Import "encoding/json" }}
            _b, _err := io.ReadAll(_httpResp.Body)
            if _err != nil {
                return nil, nil, _err
            }
            if _err = json.Unmarshal(_b, &_respObject.Body); _err != nil {
                return nil, nil, _err
            }
            {{- range $hname, $header := $resp.Headers}}
            if hdr := _httpResp.Header.Get(`{{$hname}}`); len(hdr) != 0 {
                _respObject.Header{{$hname | replace "-" "" | title}}{{if $header.Required -}} {{/*space*/}} = hdr{{- else -}}.Set(hdr){{- end}}
            }
            {{- end}}
            _resp = _respObject
        {{- else if eq $rkind "empty" -}}
            {{- $statusName := camelcase (httpStatus (atoi $code))}}
            _resp = HTTPStatus{{$statusName}}{{"{}"}}
        {{- else }}
            {{- $schema := index $resp.Content "application/json"}}
            var _respObject {{if $schema.Schema.Ref}}{{refName $schema.Schema.Ref}}{{else}}{{title $opname}}{{title $code}}Inline{{end}}
            {{- $.Import "io"}}
            {{- $.Import "encoding/json" }}
            _b, _err := io.ReadAll(_httpResp.Body)
            if _err != nil {
                return nil, nil, _err
            }
            if _err = json.Unmarshal(_b, &_respObject); _err != nil {
                return nil, nil, _err
            }
            _resp = _respObject
        {{- end -}}
    {{- end}}
    {{- if not $hasDefault}}
    default:
        {{ $.Import "fmt" -}}
        return nil, nil, fmt.Errorf("unknown response code")
    {{- end -}}
    }

    return _resp, _httpResp, nil
}
	{{end -}}
{{- end}}
