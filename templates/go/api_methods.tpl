{{range $url, $path := $.Spec.Paths -}}
    {{range $method, $op := $path.Operations -}}
        {{- $opname := lower (camelcase $op.OperationID) -}}
        {{- $.Import "net/http"}}
// {{$opname}} {{$method}} {{$url}}
func (o {{$.Name}}) {{$opname}}Op(w http.ResponseWriter, r *http.Request) error {
    var err error
    var ers []error
    var errs map[string][]string
    _, _, _ = err, ers, errs

        {{- /* Process parameters */ -}}
        {{- range $i, $param := $op.Parameters -}}
            {{- $prim := (primitive $ $param.Schema.Schema) -}}
            {{- $primNoDot := $prim | replace "." ""}}

    const n{{$i}} = `{{$param.Name}}`
    s{{$i}} {{/* This holds the space to the left */}}
            {{- if eq "query" $param.In -}}
                := r.URL.Query().Get(n{{$i}})
            {{- else if eq "header" $param.In -}}
                := r.Header.Get(n{{$i}})
            {{- else if eq "path" $param.In -}}
                := chi.URLParam(n{{$i}})
            {{- else if eq "cookie" $param.In -}}
            , err := r.Cookie(n{{$i}})
    if err != nil {
        return fmt.Errorf("failed to read cookie '{{$param.Name}}': %w", err)
    }
            {{- end}}
    var p{{$i}} {{$prim}}
            {{- if $param.Required -}}
    {{- /* Warning: This starts an else { block that covers a great deal of code */}}
    if len(s{{$i}}) == 0 {
        errs = support.AddErrs(errs, n{{$i}}, errors.New(`must not be empty`))
    } else {
            {{- else}}
    if len(s{{$i}}) != 0 {
            {{- end -}}
            {{- if ne $prim "string" -}}
                {{- $.Import "github.com/aarondl/oa3/support"}}
    p{{$i}}, err = support.StringTo{{camelcase $primNoDot}}(s{{$i}})
    if err != nil {
        errs = support.AddErrs(errs, n{{$i}}, errors.New(`was not in a valid format`))
    }
            {{- else}}
        p{{$i}} = s{{$i}}
            {{- end -}}
            {{- if mustValidate $param.Schema.Schema -}}
                {{- $.Import "github.com/aarondl/oa3/support"}}
    {{template "validate_field" (newData $ (printf "p%d" $i) $param.Schema.Schema)}}
    if len(ers) != 0 {
        errs = support.AddErrs(errs, n{{$i}}, ers...)
    }
            {{end -}}
            {{- if $param.Schema.Format -}}
            {{- $.Import "github.com/aarondl/oa3/support"}}
    if newErrs := support.Validate{{camelcase $param.Schema.Format}}(p{{$i}}); newErrs != nil {
        errs = support.AddErrs(errs, n{{$i}}, newErrs...)
    }
            {{end -}}
    }{{- /* This bracket closes the validation if above */ -}}
        {{- end}}
        {{- if $op.RequestBody -}}
            {{- $json := index $op.RequestBody.Content "application/json"}}
    // {{$json.Schema.Ref}}
    var rb {{if $json.Schema.Nullable}}*{{end}}{{refName $json.Schema.Ref}}

            {{if $op.RequestBody.Required -}}
    if r.Body == nil {
        return support.ErrNilBody
    } else {
            {{else -}}
    if r.Body != nil {
            {{end -}}
        defer r.Body.Close()
        b, err := ioutil.ReadAll(r.Body)
        if err != nil {
            return err
        }

        if err = json.Unmarshal(b, {{if not $json.Schema.Nullable}}&{{end}}rb); err != nil {
            return err
        }

        if newErrs := rb.ValidateSchema{{$.Name}}(); newErrs != nil {
            errs = support.MergeErrs(errs, newErrs)
        }
    }
        {{end}}

    if errs != nil {
        return o.cnv(errs)
    }

    return nil
}
    {{- end -}}
{{end -}}
