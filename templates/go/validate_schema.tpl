{{- /* Top level object is a Schema, Name is the name of the local var */ -}}
{{- $.Import "github.com/aarondl/oa3/support" -}}

{{- /* Validate schema helper recursively validates schema pieces */ -}}
{{- define "validate_schema_helper" -}}
    {{- $s := $.Object.Schema -}}
    {{- if eq $s.Type "array" -}}
        {{if $.Object.MaxItems -}}
    if err := support.ValidateMaxItems(o, {{$.Object.MaxItems}}); err != nil {
        {{- $.Import "fmt" -}}
        ers = append(ers, err)
    }
        {{- end -}}
        {{- if $.Object.MinItems}}
    if err := support.ValidateMinItems({{$.Name}}, {{$.Object.MinItems}}); err != nil {
        ers = append(ers, err)
    }
        {{end -}}
    for i, {{$.Name}} := range {{$.Name}} {
        var ers []error
        {{- $.Import "fmt"}}
        ctx = append(ctx, fmt.Sprintf("[%d]", i))
        {{template "validate_schema_helper" (newData $ $.Name $s.Items)}}
        {{- $.Import "strings"}}
        errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
        ctx = ctx[:len(ctx)-1]
    }
    {{else if eq $s.Type "object" -}}
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
    for k, {{$.Name}} := range {{$.Name}} {
        var ers []error
        ctx = append(ctx, k)
            {{template "validate_schema_helper" (newData $ $.Name $s.AdditionalProperties) }}
            {{- $.Import "strings"}}
        errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
        ctx = ctx[:len(ctx)-1]
    }
        {{- else if $s.Properties -}}
            {{- /* Process regular struct fields */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and (not $element.Ref) (mustValidate $element.Schema)}}
    {{template "validate_field" (recurseData $ (printf ".%s" (camelcase $name)) $element.Schema)}}
    if len(ers) != 0 {
        ctx = append(ctx, {{printf "%q" $name}})
                {{- $.Import "strings"}}
        errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
        ctx = ctx[:len(ctx)-1]
    }
                {{- end -}}
            {{- end -}}

            {{- /* Process embedded structs */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and ($element.Ref) (mustValidate $element.Schema)}}
    if newErrs := {{$name}}.{{$element}}.Validate{{$name}}(); newErrs != nil {
        errs = support.MergeErrs(errs, newErrs)
    }
                {{- end -}}
            {{- end -}}
        {{- end}}
    {{- else -}}
        {{- if mustValidate $.Object.Schema -}}
            {{- template "validate_field" (newData $ "o" $.Object.Schema) -}}
        {{- end -}}
    {{- end}}

{{- end}}

// ValidateSchema{{$.Name}} validates the object and returns
// errors that can be returned to the user.
func (o {{$.Name}}) ValidateSchema{{$.Name}}() support.Errors {
    {{- $s := $.Object.Schema}}
    var ctx []string
    var ers []error
    var errs support.Errors
    _, _ = ers, errs

    {{template "validate_schema_helper" (newData $ "o" $.Object)}}

    errs = support.AddErrs(errs, "", ers...)

    return errs
}
