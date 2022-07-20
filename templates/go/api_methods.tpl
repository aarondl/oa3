// AlreadyHandled is an interface which an error return type can optionally
// implement to stop the generated method from responding in anyway, it will
// swallow the error and not touch the ResponseWriter if this method returns
// true.
type AlreadyHandled interface {
    AlreadyHandled() bool
}

// ErrHandled is a sentinel error that implements
// the AlreadyHandled interface which prevents the
// generated handler from firing.
type ErrHandled struct {}
// Error implements error
func (ErrHandled) Error() string { return "already handled" }
// AlreadyHandled implements AlreadyHandled
func (ErrHandled) AlreadyHandled() bool { return true }

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
            {{- $prim := (primitive $ $param.Schema.Schema $param.Required)}}

    const n{{$i}} = `{{$param.Name}}`
            {{- if eq "query" $param.In}}
    s{{$i}}, s{{$i}}Exists := r.URL.Query().Get(n{{$i}}), r.URL.Query().Has(n{{$i}})
            {{- else if eq "header" $param.In}}
    s{{$i}}, s{{$i}}Exists := r.Header.Get(n{{$i}}), len(r.Header.Values(n{{$i}})) != 0
            {{- else if eq "path" $param.In -}}
                {{- $.Import "github.com/go-chi/chi/v5"}}
    s{{$i}}, s{{$i}}Exists := chi.URLParam(r, n{{$i}}), true
            {{- else if eq "cookie" $param.In -}}
                {{- $.Import "net/http"}}
    var s{{$i}} string
    s{{$i}}Exists := true
    c{{$i}}, err := r.Cookie(n{{$i}})
    if err == http.ErrNoCookie {
        s{{$i}}Exists = false
    } err != nil {
        return fmt.Errorf("failed to read cookie '{{$param.Name}}': %w", err)
    } else if err = c{{$i}}.Valid(); err != nil {
        return fmt.Errorf("failed to validate cookie '{{$param.Name}}': %w", err)
    } else {
        s{{$i}} = c{{$i}}.Value
    }
            {{- end}}
    var p{{$i}} {{$prim}}
            {{- if $param.Required -}}
    {{- /* Warning: This starts an else { block that covers a great deal of code */}}
    if !s{{$i}}Exists || len(s{{$i}}) == 0 {
            {{- $.Import "errors"}}
        errs = support.AddErrs(errs, n{{$i}}, errors.New(`must be provided and not be empty`))
    } else {
            {{- else}}
    if s{{$i}}Exists {
            {{- end -}}
            {{/*
                There are many cases for the property type:
                * string - Just assign
                * int/uint/float/bool - Convert and assign
                * (omit|null|omitnull).Val[string] - Set string without conversion
                * (omit|null|omitnull).Val[int/uint/float/bool] - Convert then set
            */}}
            {{- if ne $prim "string" -}}
                {{- $.Import "github.com/aarondl/oa3/support" -}}
                {{- $primRaw := primitiveRaw $ $param.Schema.Schema -}}
                {{- if eq $primRaw "string"}}
        p{{$i}}.Set(s{{$i}})
        err = nil
                {{- else -}}
                    {{- $convFn := printf "support.StringToBool(s%d)" $i -}}
                    {{- if eq $primRaw "chrono.DateTime" -}}
                        {{- $convFn = printf "support.StringToChronoDateTime(s%d)" $i -}}
                    {{- else if eq $primRaw "chrono.Date" -}}
                        {{- $convFn = printf "support.StringToChronoDate(s%d)" $i -}}
                    {{- else if eq $primRaw "chrono.Time" -}}
                        {{- $convFn = printf "support.StringToChronoTime(s%d)" $i -}}
                    {{- else if eq $primRaw "time.Time" -}}
                        {{- $primFmt := printf "%s" $param.Schema.Schema.Format -}}
                        {{- if eq $primFmt "date-time" -}}
                            {{- $convFn = printf "support.StringToDateTime(s%d)" $i -}}
                        {{- else if eq $primFmt "date" -}}
                            {{- $convFn = printf "support.StringToDate(s%d)" $i -}}
                        {{- else if eq $primFmt "time" -}}
                            {{- $convFn = printf "support.StringToTime(s%d)" $i -}}
                        {{- end -}}
                    {{- else if eq $primRaw "time.Duration" -}}
                        {{- $convFn = printf "support.StringToDuration(s%d)" $i -}}
                    {{- else if hasPrefix "int" $primRaw -}}
                        {{- $convFn = printf "support.StringToInt[%s](s%d, %s)" $primRaw $i (primitiveBits $ $param.Schema.Schema) -}}
                    {{- else if hasPrefix "uint" $primRaw -}}
                        {{- $convFn = printf "support.StringToUint[%s](s%d, %s)" $primRaw $i (primitiveBits $ $param.Schema.Schema) -}}
                    {{- else if hasPrefix "float" $primRaw -}}
                        {{- $convFn = printf "support.StringToFloat[%s](s%d, %s)" $primRaw $i (primitiveBits $ $param.Schema.Schema) -}}
                    {{- end -}}
                    {{- if or (hasPrefix "null." $prim) (hasPrefix "omit." $prim) (hasPrefix "omitnull." $prim)}}
        p{{$i}}c, err := {{$convFn}}
        p{{$i}}.Set(p{{$i}}c)
                    {{- else}}
        p{{$i}}, err = {{$convFn}}
                    {{- end -}}
                {{- end}}
        if err != nil {
                {{- $.Import "errors"}}
            errs = support.AddErrs(errs, n{{$i}}, errors.New(`was not in a valid format`))
        }
            {{- else}}
        p{{$i}} = s{{$i}}
            {{- end -}}
            {{- if mustValidate $param.Schema.Schema -}}
                {{- $.Import "github.com/aarondl/oa3/support"}}
        {{template "validate_field" (newDataRequired $ (printf "p%d" $i) $param.Schema.Schema $param.Required)}}
        if len(ers) != 0 {
            errs = support.AddErrs(errs, n{{$i}}, ers...)
        }
            {{end -}}
    }{{- /* This bracket closes the validation if above */ -}}
        {{- end}}
        {{$json := ""}}
        {{- if $op.RequestBody -}}
            {{- $json = index $op.RequestBody.Content "application/json"}}

    var reqBody{{" " -}}
            {{- if $json.Schema.Ref -}}
                {{- if $json.Schema.Nullable}}*{{end -}}
                {{- refName $json.Schema.Ref -}}
            {{- else -}}
                {{title $op.OperationID}}Inline
            {{- end}}

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

        if newErrs := Validate(reqBody); newErrs != nil {
            errs = support.MergeErrs(errs, newErrs)
        }
    }

        {{end}}
    if errs != nil {
        return o.converter(errs)
    }

    ret, err := o.impl.{{title $op.OperationID}}(w, r
        {{- if $op.RequestBody -}},{{" " -}}
            {{- if and $json.Schema.Ref (not $json.Schema.Nullable) (not (isInlinePrimitive $json.Schema.Schema)) -}}&{{- end -}}
            {{- if and (isInlinePrimitive $json.Schema.Schema) (not (eq $json.Schema.Schema.Type "object")) (not (eq $json.Schema.Schema.Type "array")) -}}
                {{- $p := primitive $ $json.Schema.Schema $op.RequestBody.Required}}{{$p}}(reqBody)
            {{- else -}}reqBody{{- end -}}
        {{- end -}}
        {{- range $i, $param := $op.Parameters -}}
        , p{{$i}}
        {{- end -}}
    )
    if err != nil {
        if alreadyHandled, ok := err.(AlreadyHandled); ok {
            if alreadyHandled.AlreadyHandled() {
                return nil
            }
        }
        return err
    }

    switch respBody := ret.(type) {
    {{- range $code, $resp := $op.Responses}}
    {{- $rkind := responseKind $op $code}}
    case{{" " -}}
        {{- if eq $rkind "wrapped" -}}
            {{title $op.OperationID}}{{$code}}WrappedResponse:
            {{if gt (len $resp.Headers) 0 -}}
            headers := w.Header()
                {{- range $hname, $header := $resp.Headers -}}
                    {{- $headername := $hname | replace "-" "" | title -}}
                    {{- if not $header.Required}}
            if val, ok := respBody.Header{{$headername}}.Get(); ok {
                headers.Set("{{$hname}}", val)
            }
                    {{- else}}
            headers.Set("{{$hname}}", respBody.Header{{$headername}})
                    {{- end -}}
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
        {{- else if eq $rkind "empty" -}}
            {{- $statusName := camelcase (httpStatus (atoi $code)) -}}
            HTTPStatus{{$statusName}}:
                w.WriteHeader({{$code}})
        {{- else -}}
            {{- $schema := (index $resp.Content "application/json").Schema -}}
            {{- if $schema.Ref -}}
                {{- if $schema.Nullable}}*{{end}}{{- refName $schema.Ref -}}:
            {{- else -}}
                {{title $op.OperationID}}{{title $code}}Inline:
            {{- end -}}
            {{- if ne $code "default"}}
            w.WriteHeader({{$code}})
            {{end -}}
            {{- $.Import "github.com/aarondl/oa3/support"}}
            if err := support.WriteJSON(w, respBody); err != nil {
                return err
            }
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
