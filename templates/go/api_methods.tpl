// AlreadyHandled is an interface which an error return type can optionally
// implement to stop the generated method from responding in anyway, it will
// swallow the error and not touch the ResponseWriter if this method returns
// true.
type AlreadyHandled interface {
    AlreadyHandled() bool
}

// ErrHandled is a sentinel error that implements
// the AlreadyHandled interface which prevents the
// generated handler's response code from firing.
type ErrHandled struct {}
// Error implements error
func (ErrHandled) Error() string { return "already handled" }
// AlreadyHandled implements AlreadyHandled
func (ErrHandled) AlreadyHandled() bool { return true }

{{- range $url, $path := $.Spec.Paths -}}
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
        {{- $queryOnce := false -}}
        {{- $cookieOnce := false -}}
        {{- range $i, $param := $op.Parameters -}}
            {{- $ptype := paramTypeName $ $op.OperationID $method $param}}
    const n{{$i}} = `{{$param.Name}}`
            {{- if eq "query" $param.In -}}
                {{- if not $queryOnce}}
    query := r.URL.Query()
                    {{- $queryOnce = true -}}
                {{- end}}
    s{{$i}} := query[n{{$i}}]
    s{{$i}}Exists := len(s{{$i}}) > 0 && len(s{{$i}}[0]) > 0
            {{- else if eq "header" $param.In}}
    s{{$i}} := r.Header[http.CanonicalHeaderKey(n{{$i}})]
    s{{$i}}Exists := len(s{{$i}}) > 0 && len(s{{$i}}[0]) > 0
            {{- else if eq "path" $param.In -}}
                {{- $.Import "github.com/go-chi/chi/v5"}}
    s{{$i}}, s{{$i}}Exists := []string{chi.URLParam(r, n{{$i}})}, true
            {{- else if eq "cookie" $param.In -}}
                {{- $.Import "net/http"}}
                {{- if not $cookieOnce}}
    cookies := r.Cookies()
                    {{- $cookieOnce = true -}}
                {{- end -}}
    var s{{$i}} []string
    s{{$i}}Exists := false
    for _, c := range cookies {
        if c.Name == n{{$i}} {
            if err := c.Valid(); err != nil {
                return fmt.Errorf("failed to validate cookie '{{$param.Name}}': %w", err)
            }
            s{{$i}} = append(s{{$i}}, c.Value)
            s{{$i}}Exists = s{{$i}}Exists || len(c.Value) > 0
        }
    }
            {{- end}}
    var p{{$i}} {{omitnullWrap $ $ptype $param.Schema.Nullable $param.Required }}
            {{- if $param.Required}}
    if !s{{$i}}Exists || len(s{{$i}}) == 0 {
            {{- $.Import "errors"}}
        errs = support.AddErrs(errs, n{{$i}}, errors.New(`must be provided and not be empty`))
    } else {
            {{- else}}
    if s{{$i}}Exists {
            {{- end -}}

            {{- $setVar := printf "s%d[0]" $i -}}
            {{- $mustConvert := or (ne $param.Schema.Type "string") (and (not $param.Schema.Enum) (ne $ptype "string"))}}
            {{- if $mustConvert }}
                {{- $setVar = printf "c%d" $i -}}
        c{{$i}}, err := {{paramConvertFn $ $param $ptype (printf "s%d" $i)}}
        if err != nil {
            {{- $.Import "fmt"}}
            return fmt.Errorf("failed to convert parameter %q to %q: %w", n{{$i}}, `{{$ptype}}`, err)
        }
            {{- end -}}

            {{- if $param.Schema.Enum -}}
                {{- if $mustConvert -}}
                    {{- $setVar = printf "%s(c%d)" $ptype $i -}}
                {{- else -}}
                    {{- $setVar = printf "%s(s%d[0])" $ptype $i -}}
                {{- end -}}
            {{- end -}}

            {{- if omitnullIsWrapped $param.Schema.Schema.Nullable $param.Required}}
        p{{$i}}.Set({{$setVar}})
            {{- else}}
        p{{$i}} = {{$setVar}}
            {{- end -}}

            {{- /* Validation */ -}}
            {{- if mustValidate $ $param.Schema.Schema -}}
                {{- if paramRequiresType $param }}
        if newErrs := Validate({{paramSchemaName $op.OperationID $method $param.Name}}({{omitnullUnwrap (printf "p%d" $i) $param.Schema.Nullable $param.Required}})); newErrs != nil {
            errs = support.AddErrsFlatten(errs, n{{$i}}, newErrs)
        }
                {{- else -}}
                {{- $.Import "github.com/aarondl/oa3/support" -}}
        {{- template "validate_field" (newDataRequired $ (omitnullUnwrap (printf "p%d" $i) $param.Schema.Nullable $param.Required) $param.Schema.Schema $param.Required)}}
        if len(ers) != 0 {
            errs = support.AddErrs(errs, n{{$i}}, ers...)
        }
                {{- end -}}
            {{- end}}

    }{{- /* This bracket closes the exists check if above */ -}}
        {{- end}}
        {{$json := ""}}
        {{- if $op.RequestBody -}}
            {{- $json = index $op.RequestBody.Content "application/json" -}}
            {{- if $json}}

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
            {{- end -}}
        {{- end}}
    if errs != nil {
        return o.converter(errs)
    }

    ret, err := o.impl.{{title $op.OperationID}}(w, r
        {{- if and $op.RequestBody $json -}},{{" " -}}
            {{- if and $json.Schema.Ref (not $json.Schema.Nullable) (not (isInlinePrimitive $json.Schema.Schema)) -}}&{{- end -}}
            {{- if and (isInlinePrimitive $json.Schema.Schema) (not (eq $json.Schema.Schema.Type "object")) (not (eq $json.Schema.Schema.Type "array")) -}}
                {{- $p := primitive $ $json.Schema.Schema}}{{$p}}(reqBody)
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

        {{$multi := gt (len $op.Responses) 1 -}}
        {{- $respVar := "ret" -}}
        {{- if $multi}}
    switch respBody := ret.(type) {
            {{$respVar = "respBody"}}
        {{- else}}
        _ = ret
        {{- end -}}
        {{- range $code, $resp := $op.Responses}}
            {{- $rtypeName := responseTypeName $op $code false -}}
            {{- $wrapped := responseNeedsWrap $op $code -}}
            {{- if $multi}}
    case{{" " -}}
            {{- end -}}
            {{if $wrapped -}}
                {{if $multi}}{{$rtypeName}}:{{end}}
                {{if gt (len $resp.Headers) 0 -}}
            headers := w.Header()
                    {{- range $hname, $header := $resp.Headers -}}
                        {{- $headername := $hname | replace "-" "" | title -}}
                        {{- if not $header.Required}}
            if val, ok := {{$respVar}}.Header{{$headername}}.Get(); ok {
                headers.Set("{{$hname}}", val)
            }
                        {{- else}}
            headers.Set("{{$hname}}", {{$respVar}}.Header{{$headername}})
                        {{- end -}}
                    {{- end -}}
                {{- end -}}
                {{- if ne $code "default"}}
            w.WriteHeader({{$code}})
                {{- end -}}
                {{- if $resp.Content -}}
                    {{- $.Import "github.com/aarondl/oa3/support"}}
            if err := support.WriteJSON(w, {{$respVar}}); err != nil {
                return err
            }
                {{- end -}}
            {{- else if not $resp.Content -}}
                {{- $statusName := camelcase (httpStatus (atoi $code)) -}}
                {{if $multi}}HTTPStatus{{$statusName}}:{{end}}
                w.WriteHeader({{$code}})
            {{- else -}}
                {{- if $multi -}}
                    {{- with $schema := index $resp.Content "application/json" -}}
                        {{- if $schema.Schema.Ref -}}
                            {{- if $schema.Schema.Nullable}}*{{end}}{{- refName $schema.Schema.Ref -}}:
                        {{- else -}}
                            {{title $op.OperationID}}{{title $code}}Inline:
                        {{- end -}}
                    {{- else -}}
                            {{- refName $schema.Ref -}}:
                    {{- end -}}
                {{- end -}}
                {{- if ne $code "default"}}
            w.WriteHeader({{$code}})
                {{end -}}
                {{- with $schema := index $resp.Content "application/json" -}}
                    {{- $.Import "github.com/aarondl/oa3/support"}}
            if err := support.WriteJSON(w, {{$respVar}}); err != nil {
                return err
            }
                {{- else -}}
            if {{$respVar}} != nil {
                {{- $.Import "io"}}
                if _, err := io.Copy(w, {{$respVar}}); err != nil {
                    return err
                }
                if err := {{$respVar}}.Close(); err != nil {
                    return err
                }
            }
                {{- end -}}
        {{- end -}}
    {{- end -}}
    {{- if $multi}}
    default:
        _ = respBody
        panic("impossible case")
    }
    {{- end}}

    return nil
}
    {{- end -}}
{{end -}}
