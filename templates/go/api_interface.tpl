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
    // {{wrapWith 70 "\n// " (trimSuffix "\n" (replace "\n" "" $op.Description))}}
        {{end -}}
    {{$opname}}(w http.ResponseWriter, r *http.Request
        {{- if $op.RequestBody -}}
            {{- $json := index $op.RequestBody.Content "application/json" -}}
            {{- if $json -}}
            , body{{" " -}}
                {{- if $json.Schema.Ref -}}
                    {{- if not (isInlinePrimitive $json.Schema.Schema) -}}*{{- end -}}
                    {{- title (refName $json.Schema.Ref) -}}
                {{- else if not (or (eq $json.Schema.Schema.Type "object") (eq $json.Schema.Schema.Type "array")) -}}
                    {{- primitiveWrapped $ $json.Schema.Schema $json.Schema.Nullable $op.RequestBody.Required -}}
                {{- else -}}
                    {{title $op.OperationID}}Inline
                {{- end -}}
            {{- end -}}
        {{- end -}}
        {{- range $param := $op.Parameters -}}
        , {{untitle (camelcase $param.Name)}} {{omitnullWrap $ (paramTypeName $ $op.OperationID $method $param) $param.Schema.Nullable $param.Required }}
        {{- end -}}
        {{- $respName := responseRefName $op -}}
        {{- if $respName -}}
    ) ({{if responseNeedsPtr $op}}*{{end}}{{$respName}}, error)
        {{else -}}
    ) error
        {{end -}}
    {{- end -}}
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
func Validate(toValidate validatable) support.Errors {
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
            {{- $json := index $op.RequestBody.Content "application/json" -}}
            {{- if $json -}}
            , body{{" " -}}
                {{- if $json.Schema.Ref -}}
                    {{- if not (isInlinePrimitive $json.Schema.Schema) -}}*{{- end -}}
                    {{$.Params.package}}.{{- refName $json.Schema.Ref -}}
                {{- else if not (or (eq $json.Schema.Schema.Type "object") (eq $json.Schema.Schema.Type "array")) -}}
                    {{- primitiveWrapped $ $json.Schema.Schema $json.Schema.Nullable $op.RequestBody.Required -}}
                {{- else -}}
                    {{$.Params.package}}.{{title $op.OperationID}}Inline
                {{- end -}}
            {{- end -}}
        {{- end -}}
        {{- range $param := $op.Parameters -}}
        , {{untitle (camelcase $param.Name)}} {{omitnullWrap $ (paramTypeName $ $op.OperationID $method $param) $param.Schema.Nullable $param.Required }}
        {{- end -}}
        {{- $respName := responseRefName $op -}}
        {{- if $respName -}}
    ) ({{if responseNeedsPtr $op}}*{{end}}{{$.Params.package}}.{{$respName}}, error) {
        {{- else -}}
    ) error {
        {{- end}}
    panic("not implemented")
}
        {{- end -}}
{{- end}}
*/
