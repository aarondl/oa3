{{- $.Import "net/http" -}}
// {{$.Name}}API is the interface that an application server must implement
// in order to use this package.
{{- if $.Spec.Info.Description}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" $.Spec.Info.Description)}}
{{end -}}
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
        r.Method(http.Method{{title (lower $op.Method)}}, `{{$op.Path}}`, eh.Wrap(o.{{$op.Op.OperationID}}Op))
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
    {{range $method, $op := $path.Operations -}}
    {{- $opname := title $op.OperationID -}}
// {{$opname}}Response one-of enforcer
//
// Implementors:
    {{- range $code, $resp := $op.Responses}}
    // - {{if $resp.Content -}}
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
    {{if $resp.Content -}}
        {{- $schema := index $resp.Content "application/json" -}}
// {{$opname}}Impl implements {{$opname}}Response({{$code}}) for {{refName $schema.Schema.Ref}}
func ({{refName $schema.Schema.Ref}}) {{$opname}}Impl() {}
        {{- else -}}
{{- $statusName := camelcase (httpStatus (atoi $code)) -}}
// {{$opname}}Impl implements {{$opname}}Response({{$code}}) for HTTPStatus{{$statusName}}
func (HTTPStatus{{$statusName}}) {{$opname}}Impl() {}
        {{- end -}}
    {{- end}}

    {{end -}}
{{end -}}

{{range $status, $_ := $needHTTPStatuses -}}
// HTTPStatus{{$status}} is an empty response
type HTTPStatus{{$status}} struct {}
{{end -}}
