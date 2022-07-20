{{- $.Import "net/http" -}}
// Interface is the interface that an application server must implement
// in order to use this package.
{{- if $.Spec.Info.Description}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" $.Spec.Info.Description)}}
{{- end}}
type Interface interface {
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
    ) ({{title $op.OperationID}}Response, error)
    {{end -}}
{{end -}}
}

{{.Import "github.com/aarondl/oa3/support" -}}
{{.Import "github.com/go-chi/chi/v5" -}}
type (
    // {{$.Name}} implements all the wrapper functionality for the API
    {{$.Name}} struct {
        impl      Interface
        converter support.ValidationConverter
        router    *chi.Mux
    }
)

// New{{$.Name}} constructor
func New{{$.Name}}(
    impl Interface,
    cnv support.ValidationConverter,
    eh support.ErrorHandler,
    mw support.MW,
    ) http.Handler {

    o := {{.Name}}{
        impl:      impl,
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


type validatable interface {
    validateSchema() support.Errors
}
// Validate a schema
func Validate[T validatable](toValidate T) support.Errors {
    return toValidate.validateSchema()
}

{{template "responses" $}}

/*
Here is a copy pastable list of function signatures
for implementing the main interface
{{range $url, $path := $.Spec.Paths -}}
    {{- range $method, $op := $path.Operations -}}
    {{- $opname := title $op.OperationID}}
// {{$opname}} {{$method}} {{$url}}
{{if $op.Description -}}
// {{wrapWith 70 "\n// " (trimSuffix "\n" $op.Description)}}
{{end -}}
func (a API) {{$opname}}(w http.ResponseWriter, r *http.Request
        {{- if $op.RequestBody -}}
        , {{$media := index $op.RequestBody.Content "application/json" -}}
        body{{" " -}}
            {{- if $media.Schema.Ref -}}
                {{- if not (isInlinePrimitive $media.Schema.Schema) -}}*{{- end -}}
                {{$.Params.package}}.{{- refName $media.Schema.Ref -}}
            {{- else if isInlinePrimitive $media.Schema.Schema -}}
                {{- primitive $ $media.Schema.Schema $op.RequestBody.Required -}}
            {{- else -}}
                {{$.Params.package}}.{{title $op.OperationID}}Inline
            {{- end -}}
        {{- end -}}
        {{- range $param := $op.Parameters -}}
        , {{untitle (camelcase $param.Name)}} {{primitive $ $param.Schema.Schema $param.Required -}}
        {{- end -}}
    ) ({{$.Params.package}}.{{title $op.OperationID}}Response, error) {
    panic("not implemented")
}
        {{- end -}}
{{- end}}
*/
