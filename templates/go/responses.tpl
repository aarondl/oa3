{{- $needHTTPStatuses := dict -}}
{{- $needIOReadClosers := dict -}}

{{- range $url, $path := $.Spec.Paths -}}
    {{- range $method, $op := $path.Operations -}}
        {{- $multi := gt (len $op.Responses) 1 -}}
        {{- $opname := title $op.OperationID}}
        {{- if $multi}}

// {{$opname}}Response one-of enforcer
//
// Implementors:
            {{- $implementorsDuplicates := dict -}}
            {{- range $code, $resp := $op.Responses}}
                {{- $typeName := responseTypeName $op $code true -}}
                {{- if not (hasKey $implementorsDuplicates $typeName) -}}
                    {{- $_ := set $implementorsDuplicates $typeName ""}}
// - {{responseTypeName $op $code false}}
                {{- end -}}
            {{- end}}
type {{$opname}}Response interface {
    {{$opname}}Impl()
}
        {{- end -}}

        {{- $duplicateResponseObjects := dict -}}
        {{- range $code, $resp := $op.Responses}}
            {{- $typeName := responseTypeName $op $code true -}}
            {{- if responseNeedsWrap $op $code -}}
                {{- if not (hasKey $duplicateResponseObjects $typeName) }}
                    {{- $_ := set $duplicateResponseObjects $typeName "" -}}
                    {{- $wrapName := responseTypeName $op $code false}}
// {{$wrapName}} wraps the normal body response with a
// struct to be able to additionally return headers or differentiate between
// multiple response codes with the same response body.
type {{$wrapName}} struct {
                    {{- if responseNeedsCodeWrap $op $code}}
    Code int
                    {{- end -}}
                    {{- if responseNeedsHeaderWrap $op $code -}}
                        {{- range $hname, $header := $resp.Headers}}
    Header{{$hname | replace "-" "" | title}} {{if $header.Required -}}
                                    string
                                {{- else -}}
                                    {{- $.Import "github.com/aarondl/opt/omit" -}}
                                    omit.Val[string]
                                {{- end -}}
                        {{- end}}
                        {{- with $content := $resp.Content -}}
                        {{- if not (index $content "application/json")}}
                            {{- $_ := set $needIOReadClosers $typeName "" -}}
                            {{- end -}}
                        {{- end -}}
                    {{- end}}
    Body {{$typeName}}
}

                    {{if $multi}}
// {{$opname}}Impl implements {{$opname}}Response({{$code}}) for {{$wrapName}}
func ({{$wrapName}}) {{$opname}}Impl() {}
                    {{- end -}}
                {{- end -}}

            {{- else if not $resp.Content -}}
                {{- $statusName := camelcase (httpStatus (atoi $code))}}
                {{- $_ := set $needHTTPStatuses $statusName "" -}}
                {{- if $multi}}
// {{$opname}}Impl implements {{$opname}}Response({{$code}}) for {{$typeName}}
func ({{$typeName}}) {{$opname}}Impl() {}
                {{- end -}}
            {{- else -}}
                {{- if not (index $resp.Content "application/json")}}
                    {{- $_ := set $needIOReadClosers $typeName "" -}}
                {{- end -}}
                {{- if $multi -}}
                    {{- range $content, $schema := $resp.Content}}
// {{$opname}}Impl implements {{$opname}}Response({{$code}}) for {{$typeName}}
func ({{$typeName}}) {{$opname}}Impl() {}
                    {{- end -}}
                {{- end -}}
            {{- end -}}
        {{- end -}}
    {{- end -}}
{{- end}}

{{range $status, $_ := $needHTTPStatuses -}}
// HTTPStatus{{$status}} is an empty response
type HTTPStatus{{$status}} struct {}
{{end -}}

{{range $name, $_ := $needIOReadClosers -}}
    {{- $.Import "io" -}}
type {{$name}} io.ReadCloser
{{end -}}
