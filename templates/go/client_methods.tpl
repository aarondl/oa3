{{- range $url, $path := $.Spec.Paths -}}
    {{- range $method, $op := $path.Operations -}}
        {{- $opname := (title $op.OperationID) -}}
        {{- $respRefName := responseRefName $op -}}
        {{- $.Import "net/http"}}
// {{$opname}} {{$method}} {{$url}}
{{- if $op.Description}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" (replace "\n" " " $op.Description))}}
{{- end -}}
{{- $.Import "context"}}
func (_c Client) {{$opname}}(ctx context.Context
        {{- if hasComplexServers $op.Servers -}}
        , baseURL URLBuilder{{$url | filterNonIdentChars | title}}{{$method | filterNonIdentChars | title}}
        {{- else if hasComplexServers $path.Servers -}}
        , baseURL URLBuilder{{$url | filterNonIdentChars | title}}
        {{- end -}}
		{{- if $op.RequestBody -}}
            {{- with $json := index $op.RequestBody.Content "application/json" -}}
                , body{{" " -}}
                    {{- if $json.Schema.Ref -}}
                        {{- if not (isInlinePrimitive $json.Schema.Schema) -}}*{{- end -}}
                        {{- title (refName $json.Schema.Ref) -}}
                    {{- else if not (or (eq $json.Schema.Schema.Type "object") (eq $json.Schema.Schema.Type "array")) -}}
                        {{- primitiveWrapped $ $json.Schema.Schema $json.Schema.Nullable $op.RequestBody.Required -}}
                    {{- else -}}
                        {{title $op.OperationID}}Inline
                    {{- end -}}
            {{- else -}}
                {{- $.Import "io" -}}
                , body io.ReadCloser
            {{- end -}}
        {{- end -}}
		{{- range $param := $op.Parameters -}}
		, {{untitle (camelcase $param.Name)}} {{omitnullWrap $ (paramTypeName $ $op.OperationID $method $param) $param.Schema.Nullable $param.Required}}
		{{- end -}}
	) ({{$respRefName}}, *http.Response, error) {
        var _resp {{$respRefName}}
        var _httpResp *http.Response
        var _err error
    {{- if and (not $op.Servers) (not $path.Servers)}}
    baseURL := _c.url
    {{- else if and $op.Servers (not (hasComplexServers $op.Servers))}}
    baseURL := {{(index $op.Servers 0).URL | filterNonIdentChars | title}}
    {{- else if and $path.Servers (not (hasComplexServers $path.Servers))}}
    baseURL := {{(index $path.Servers 0).URL | filterNonIdentChars | title}}
    {{- end -}}
    {{- $.Import "strings"}}
    _urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `{{$url}}`
	{{- range $param := $op.Parameters -}}
        {{- if and (eq $param.In "path") -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- $.Import "strings"}}
            _urlStr = strings.Replace(_urlStr, `{{"{"}}{{$param.Name}}{{"}"}}`, fmt.Sprintf("%v", {{$pname}}), 1) 
        {{- end -}}
	{{- end}}
	_req, _err := http.NewRequestWithContext(ctx, http.Method{{camelcase $method}}, _urlStr, nil)
    if _err != nil {
        return _resp, _httpResp, _err
    }
	{{- if $op.RequestBody -}}
        {{- $json := index $op.RequestBody.Content "application/json" -}}
        {{- if $json -}}
            {{- $.Import "github.com/aarondl/json" -}}
            {{- $.Import "bytes" -}}
            {{- $.Import "io"}}
            _bodyBytes, _err := json.Marshal(body)
            if _err != nil {
                return _resp, _httpResp, _err
            }
            _req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
        {{- else}}
            _req.Body = body
        {{- end -}}
    {{- end -}}
	{{- range $param := $op.Parameters -}}
        {{- if and (eq $param.In "header") (not (has (lower $param.Name) (list "accept" "content-type" "authorization"))) -}}
            {{- $pname := untitle (camelcase $param.Name) -}}
            {{- $wrapped := omitnullIsWrapped $param.Schema.Schema.Nullable $param.Required -}}
            {{- $useName := $pname -}}
            {{- if $wrapped -}}
                {{- $useName = "_val"}}
    if _val, _ok := {{$pname}}.Get(); _ok {
            {{- end -}}
            {{- if eq $param.Schema.Schema.Type "array" -}}
                {{- if deref $param.Explode}}
        for _, _v := range {{$useName}} {
            _req.Header.Add(`{{$param.Name}}`, fmt.Sprintf("%v", _v))
        }
                {{- else}}
        var _{{$pname}}Slice []string
        for _, _v := range {{$useName}} {
            _{{$pname}}Slice = append(_{{$pname}}Slice, fmt.Sprintf("%v", _v))
        }
                {{- $.Import "strings" }}
        _req.Header.Set(`{{$param.Name}}`, strings.Join(_{{$pname}}Slice, ","))
                {{- end -}}
            {{- else}}
        _req.Header.Add(`{{$param.Name}}`, fmt.Sprintf("%v", {{$useName}}))
            {{- end -}}
            {{if $wrapped -}}
    }
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
            {{- $wrapped := omitnullIsWrapped $param.Schema.Schema.Nullable $param.Required -}}
            {{- $useName := $pname -}}
            {{- if $wrapped -}}
                {{- $useName = "_val"}}
    if _val, _ok := {{$pname}}.Get(); _ok {
            {{- end -}}
            {{- if eq $param.Schema.Schema.Type "array" -}}
                {{- if deref $param.Explode}}
        for _, _v := range {{$useName}} {
            _query.Add(`{{$param.Name}}`, fmt.Sprintf("%v", _v))
        }
                {{- else}}
        var _{{$pname}}Slice []string
        for _, _v := range {{$useName}} {
            _{{$pname}}Slice = append(_{{$pname}}Slice, fmt.Sprintf("%v", _v))
        }
                {{- $.Import "strings" }}
        _query.Set(`{{$param.Name}}`, strings.Join(_{{$pname}}Slice, ","))
                {{- end -}}
            {{- else}}
        _query.Add(`{{$param.Name}}`, fmt.Sprintf("%v", {{$useName}}))
            {{- end -}}
            {{if $wrapped -}}
    }
            {{- end -}}
        {{- end -}}
	{{- end}}
    if len(_query) > 0 {
        _req.URL.RawQuery = _query.Encode()
    }

	{{range $param := $op.Parameters -}}
        {{- if eq $param.In "cookie" -}}
    // $param.Name cookie param not supported by go client yet
            {{fail "cookie param not supported by go client yet"}}
        {{- end -}}
	{{- end -}}

    _httpResp, _err = _c.doRequest(ctx, _req)
    if _err != nil {
        return _resp, _httpResp, _err
    }

    switch _httpResp.StatusCode {
    {{- $hasDefault := index $op.Responses "default" -}}
    {{- $multi := gt (len $op.Responses) 1}}
    {{- range $code, $resp := $op.Responses}}
        {{- $typeName := responseTypeName $op $code true}}
        {{if not (eq $code "default")}}case {{$code}}:{{- else -}}default:{{- end}}
        {{- if responseNeedsWrap $op $code -}}
            {{- $wrapName := responseTypeName $op $code false}}
            var _respObject {{$wrapName}}
            {{- with $content := $resp.Content -}}
                {{- if index $content "application/json" -}}
                    {{- $.Import "io" -}}
                    {{- $.Import "github.com/aarondl/json" }}
            _b, _err := io.ReadAll(_httpResp.Body)
            if _err != nil {
                return _resp, _httpResp, _err
            }
            if _err = json.Unmarshal(_b, &_respObject.Body); _err != nil {
                return _resp, _httpResp, _err
            }
                {{- else -}}
            _respObject.Body = _httpResp.Body
                {{- end -}}
            {{- end -}}
            {{- range $hname, $header := $resp.Headers}}
            if hdr := _httpResp.Header.Get(`{{$hname}}`); len(hdr) != 0 {
                _respObject.Header{{$hname | replace "-" "" | title}}{{if $header.Required -}} {{/*space*/}} = hdr{{- else -}}.Set(hdr){{- end}}
            }
            {{- end}}
            _resp = _respObject
        {{- else if not $resp.Content -}}
            _resp = {{$typeName}}{{"{}"}}
        {{- else -}}
            {{- with $content := $resp.Content -}}
                {{- with $schema := index $content "application/json" -}}
                    {{- $.Import "io" -}}
                    {{- $.Import "github.com/aarondl/json" }}
            var _respObject {{$typeName}}
            _b, _err := io.ReadAll(_httpResp.Body)
            if _err != nil {
                return _resp, _httpResp, _err
            }
            if _err = json.Unmarshal(_b, &_respObject); _err != nil {
                return _resp, _httpResp, _err
            }
            _resp = _respObject
                {{- else -}}
            _resp = _httpResp.Body
                {{- end -}}
            {{- end -}}
        {{- end -}}
    {{- end}}
    {{- if not $hasDefault}}
    default:
        {{ $.Import "fmt" -}}
        return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
    {{- end -}}
    }

    return _resp, _httpResp, nil
}
	{{end -}}
{{- end}}
