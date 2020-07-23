{{- /* Top level object is a Schema, Name is the name of the local var */ -}}
{{- $.Import "github.com/aarondl/oa3/support" -}}

{{- /* Validate schema helper recursively validates schema pieces */ -}}
{{- define "validate_schema_helper" -}}
    {{- $s := $.Object.Schema -}}
    {{- if eq $s.Type "array" -}}
        {{if $.Object.MaxItems -}}
    if err := support.MaxItems({{$.Name}}, {{$.Object.MaxItems}}); err != nil {
        ers = append(ers, err)
    }
        {{end -}}
        {{- if $.Object.MinItems}}
    if err := support.MinItems({{$.Name}}, {{$.Object.MinItems}}); err != nil {
        ers = append(ers, err)
    }
        {{end -}}
    for _, v := range {{$.Name}} {
        {{template "validate_schema_helper" (named $ "o" $s.Items)}}
    }
    {{else if eq $s.Type "object" -}}
        {{- if $s.AdditionalProperties -}}
            {{- if not $s.AdditionalProperties.SchemaRef -}}{{fail "additionalItems being bool is not supported"}}{{- end}}
    // Min/Max Properties
    for k, v := range {{$.Name}} {
        {{template "validate_schema_helper" (named $ "o" $s.AdditionalProperties) }}
    }
        {{- else if $s.Properties -}}
            {{- /* Process regular struct fields */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and (not $element.Ref) (mustValidate $element.Schema) -}}
    {{template "validate_field" (named $ $name $element)}}
                {{- end -}}
            {{- end -}}

            {{- /* Process embedded structs */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and (not $element.Ref) (mustValidate $element.Schema)}}
                // {{$name}}.Validate{{$name}}()
                {{- end -}}
            {{- end -}}
        {{- end}}
    {{- else -}}
        {{- template "validate_field" (namedNoRecurse $ "o" $.Object) -}}
    {{- end}}

{{- end}}

// ValidateSchema{{$.Name}} validates the object and returns
// errors that can be returned to the user.
func (o {{$.Name}}) ValidateSchema{{$.Name}}() support.Errors {
    {{- $s := $.Object.Schema -}}
    var errs support.Errors

    {{template "validate_schema_helper" (namedNoRecurse $ "o" $.Object)}}

    return errs
}
