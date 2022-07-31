{{- /* Top level object is a Schema, Name is the name of the local var */ -}}
{{- $.Import "github.com/aarondl/oa3/support" -}}

{{- /* Validate schema helper recursively validates schema pieces */ -}}
{{- define "validate_schema_helper" -}}
    {{- $s := $.Object.Schema -}}
    {{- if eq $s.Type "array" -}}
        {{if $.Object.MaxItems -}}
    if err := support.ValidateMaxItems(o, {{$.Object.MaxItems}}); err != nil {
        ers = append(ers, err)
    }
        {{- end -}}
        {{- if $.Object.MinItems}}
    if err := support.ValidateMinItems({{$.Name}}, {{$.Object.MinItems}}); err != nil {
        ers = append(ers, err)
    }
        {{end -}}
        {{- if mustValidateRecurse $s.Items.Schema}}
    for i, {{$.Name}} := range {{$.Name}} {
        _ = {{$.Name}}
        {{- $.Import "fmt"}}
        ctx = append(ctx, fmt.Sprintf("[%d]", i))
            {{- if $s.Items.Ref}}
        if newErrs := Validate({{$.Name}}); newErrs != nil {
            errs = support.AddErrsFlatten(errs, strings.Join(ctx, "."), newErrs)
        }
            {{- else }}
        var ers []error
        {{template "validate_schema_helper" (newDataRequired $ $.Name $s.Items true)}}
        errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
            {{- end -}}
        {{- $.Import "strings"}}
        ctx = ctx[:len(ctx)-1]
    }
        {{- end -}}
    {{- else if eq $s.Type "object" -}}
        {{- if $s.AdditionalProperties -}}
            {{- if not $s.AdditionalProperties.SchemaRef -}}{{fail "additionalItems being bool is not supported"}}{{- end}}
            {{if $.Object.MaxProperties -}}
    if err := support.ValidateMaxProperties({{$.Name}}, {{$.Object.MaxProperties}}); err != nil {
        ers = append(ers, err)
    }
            {{- end -}}
            {{- if $.Object.MinProperties}}
    if err := support.ValidateMinProperties({{$.Name}}, {{$.Object.MinProperties}}); err != nil {
        ers = append(ers, err)
    }
            {{end -}}
            {{- if mustValidateRecurse $s.AdditionalProperties.Schema -}}
    for k, {{$.Name}} := range {{$.Name}} {
        _ = {{$.Name}}
        var ers []error
        ctx = append(ctx, k)
            {{template "validate_schema_helper" (newDataRequired $ $.Name $s.AdditionalProperties true) }}
            {{- $.Import "strings"}}
        errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
        ctx = ctx[:len(ctx)-1]
    }
            {{- end -}}
        {{- else if $s.Properties -}}
            {{- /* Process regular struct fields */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and (not $element.Ref) (mustValidate $element.Schema) -}}
                    {{- $isRequired := $s.IsRequired $name -}}
                    {{- if and $element.Enum (gt (len $element.Enum) 0) }}
        if newErrs := Validate({{$.Name}}.{{omitnullUnwrap $ $s (camelcase $name) $element.Nullable $isRequired}}); newErrs != nil {
            errs = support.AddErrsFlatten(errs, strings.Join(ctx, "."), newErrs)
        }
                    {{- else}}
    {{template "validate_field" (recurseDataSetRequired $ (printf ".%s" (omitnullUnwrap $ $s (camelcase $name) $element.Nullable $isRequired)) $element.Schema $isRequired)}}
    if len(ers) != 0 {
        ctx = append(ctx, {{printf "%q" $name}})
                {{- $.Import "strings"}}
        errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
        ctx = ctx[:len(ctx)-1]
    }
                    {{- end -}}
                {{- end -}}
            {{- end -}}

            {{- /* Process embedded structs */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and ($element.Ref) (mustValidate $element.Schema)}}
    if newErrs := Validate(o.{{camelcase $name}}); newErrs != nil {
        ctx = append(ctx, {{printf "%q" $name}})
                {{- $.Import "strings"}}
        errs = support.AddErrsFlatten(errs, strings.Join(ctx, "."), newErrs)
        ctx = ctx[:len(ctx)-1]
    }
                {{- end -}}
            {{- end -}}
        {{- end}}
    {{- else -}}
        {{- if mustValidate $.Object.Schema -}}
            {{- template "validate_field" (newDataRequired $ "o" $.Object.Schema true) -}}
        {{- end -}}
    {{- end}}

{{- end}}

// validateSchema validates the object and returns
// errors that can be returned to the user.
func (o {{title $.Name}}) validateSchema() support.Errors {
    {{- $s := $.Object.Schema}}
    var ctx []string
    var ers []error
    var errs support.Errors
    _, _, _ = ctx, ers, errs

    {{template "validate_schema_helper" (newData $ "o" $.Object)}}

    return errs
}
