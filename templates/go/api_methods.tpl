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
            {{- $.Import "errors"}}
        errs = support.AddErrs(errs, n{{$i}}, errors.New(`must not be empty`))
    } else {
            {{- else}}
    if len(s{{$i}}) != 0 {
            {{- end -}}
            {{- if ne $prim "string" -}}
                {{- $.Import "github.com/aarondl/oa3/support"}}
    p{{$i}}, err = support.StringTo{{camelcase $primNoDot}}(s{{$i}})
    if err != nil {
            {{- $.Import "errors"}}
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
        {{$json := ""}}
        {{- if $op.RequestBody -}}
            {{- $json = index $op.RequestBody.Content "application/json"}}

    var reqBody {{if $json.Schema.Nullable}}*{{end}}{{refName $json.Schema.Ref}}

            {{if $op.RequestBody.Required -}}
    if r.Body == nil {
        return support.ErrNoBody
    } else {
            {{- else -}}
    if r.Body != nil {
            {{- end -}}
            {{- $.Import "github.com/aarondl/oa3/support"}}
        if err = support.ReadJSON(r, {{if not $json.Schema.Nullable}}&{{end}}reqBody); err != nil {
            return err
        }

        if newErrs := reqBody.ValidateSchema{{refName $json.Schema.Ref}}(); newErrs != nil {
            errs = support.MergeErrs(errs, newErrs)
        }
    }

    if errs != nil {
        return o.converter(errs)
    }
        {{end}}

    ret, err := o.impl.{{title $op.OperationID}}(w, r
        {{- if $op.RequestBody -}}
        , {{if not $json.Schema.Nullable}}&{{end}}reqBody
        {{- end -}}
        {{- range $i, $param := $op.Parameters -}}
        , p{{$i}}
        {{- end -}}
    )
    if err != nil {
        return err
    }

    switch respBody := ret.(type) {
    {{- range $code, $resp := $op.Responses}}
    case{{" " -}}
        {{- if $resp.Headers -}}
            {{title $op.OperationID}}{{$code}}HeadersResponse:
            headers := w.Header()
            {{- range $hname, $header := $resp.Headers -}}
                {{- $headername := $hname | replace "-" "" | title -}}
                {{- if not $header.Required}}
            if respBody.Header{{$headername}}.Valid {
                headers.Set("{{$hname}}", respBody.Header{{$headername}}.String)
            }
                {{- else}}
            headers.Set("{{$hname}}", respBody.Header{{$headername}})
                {{- end -}}
            {{- end -}}
            {{- if ne $code "default"}}
            w.WriteHeader({{$code}})
            {{- end -}}
            {{- if $resp.Content -}}
                {{- $.Import "github.com/aarondl/oa3/support"}}
            if err := support.WriteJSON(w, respBody); err != nil {
                return err
            }
            {{- end -}}
        {{- else if $resp.Content -}}
            {{- $schema := index $resp.Content "application/json" -}}
            {{- if $schema.Schema.Nullable}}*{{end}}{{- refName $schema.Schema.Ref -}}:
            {{- if ne $code "default"}}
            w.WriteHeader({{$code}})
            {{end -}}
            {{- $.Import "github.com/aarondl/oa3/support"}}
            if err := support.WriteJSON(w, respBody); err != nil {
                return err
            }
        {{- else -}}
            {{- $statusName := camelcase (httpStatus (atoi $code)) -}}
            HTTPStatus{{$statusName}}:
                w.WriteHeader({{$code}})
        {{- end -}}
    {{- end}}
    default:
        _ = respBody
        panic("impossible case")
    }

    return nil
}
    {{- end -}}
{{end -}}