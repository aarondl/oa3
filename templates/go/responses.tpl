{{$needHTTPStatuses := dict}}

{{range $url, $path := $.Spec.Paths -}}
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
			{{- if $schema.Schema.Ref -}}
				{{- refName $schema.Schema.Ref}} - {{ $schema.Schema.Ref -}}
			{{- else -}}
				{{title $opname}}{{title $code}}Inline
			{{- end -}}
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
