{{- $.Import "net/http" -}}
// {{$.Name}}API is the interface that an application server must implement
// in order to use this package.
{{- if $.Spec.Info.Description}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" $.Spec.Info.Description)}}
{{- end}}
type {{$.Name}}API interface {
{{range $url, $path := $.Spec.Paths -}}
    {{range $method, $op := $path.Operations -}}
    {{- $opname := title $op.OperationID -}}
    // {{$opname}} {{$method}} {{$url}}
        {{if $op.Description -}}
    // {{wrapWith 70 "\n// " (trimSuffix "\n" $op.Description)}}
        {{end -}}
    {{$opname}}(w http.ResponseWriter, r *http.Request
        {{- if $op.RequestBody -}}
        , {{$media := index $op.RequestBody.Content "application/json" -}}
        body *{{refName $media.Schema.Ref}}
        {{- end -}}
        {{- range $param := $op.Parameters -}}
        , {{untitle (camelcase $param.Name)}} {{primitive $ $param.Schema.Schema -}}
        {{- end -}}
    ) ({{title $op.OperationID}}Response, error)
        {{end -}}
{{end -}}
}

{{.Import "github.com/aarondl/oa3/support" -}}
{{.Import "github.com/go-chi/chi" -}}
type (
    // {{$.Name}} implements all the wrapper functionality for the API
    {{$.Name}} struct {
        impl      {{$.Name}}API
        converter support.ValidationConverter
        router    *chi.Mux
    }
)

// New{{$.Name}} constructor
func New{{$.Name}}(
    apiInterface {{$.Name}}API,
    cnv support.ValidationConverter,
    eh support.ErrorHandler,
    mw support.MW,
    ) http.Handler {

    o := {{.Name}}{
        impl:      apiInterface,
        converter: cnv,
        router:    chi.NewRouter(),
    }

    {{range $tag := taggedPaths $.Spec -}}
    // {{if $tag.Tag}}{{$tag.Tag}} tagged{{else}}Untagged{{end}} operations
    o.router.Group(func(r chi.Router) {
        if m, ok := mw["{{$tag.Tag}}"]; ok {
            r.Use(m...)
        }
            {{- range $op := $tag.Ops}}
        r.Method(http.Method{{title (lower $op.Method)}}, `{{$op.Path}}`, eh.Wrap(o.{{lower (camelcase $op.Op.OperationID)}}Op))
            {{- end -}}
    })
    {{end}}

    return o
}

// ServeHTTP implements http.Handler
func (o {{$.Name}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    o.router.ServeHTTP(w, r)
}

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
                                    {{- $.Import "github.com/volatiletech/null/v8" -}}
                                    null.String
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
func ({{refName $schema.Schema.Ref}}) {{$opname}}Impl() {}
            {{- end -}}
        {{- end -}}
    {{- end -}}
{{- end}}

{{range $status, $_ := $needHTTPStatuses -}}
// HTTPStatus{{$status}} is an empty response
type HTTPStatus{{$status}} struct {}
{{end -}}
