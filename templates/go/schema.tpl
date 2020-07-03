{{- /* Used to output a type name, takes a TemplateData with the object set to a schema ref */ -}}
{{- define "type_name" -}}
    {{- if or .Object.Ref (isInlinePrimitive .Object.Schema) -}}
        {{- template "schema" (named $ .Name .Object) -}}
    {{- else -}}
        {{- if .Object.Schema.Nullable}}*{{end}}{{.Name}}
    {{- end -}}
{{- end -}}

{{- /* Used to output an embedded type, takes a TemplateData with the object set
to a schema ref */ -}}
{{- define "type_embedded" -}}
    {{- if and (not .Object.Ref) (not (isInlinePrimitive .Object.Schema))}}

{{template "schema_top" $ -}}
    {{- end -}}
{{- end -}}

{{- /* Write out the schema after ensuring it's not a ref */ -}}
{{- define "schema_noref" -}}
    {{- $s := .Object -}}

    {{- if $s.Enum -}}
        {{- if not (eq $s.Type "string") -}}{{fail "non-string enums not supported"}}{{- end -}}
        string

const ({{range $value := $s.Enum}}
    {{$.Name}}{{camelcase $value}} {{$.Name}} = "{{$value}}"
        {{- end}}
)

    {{- else if or $s.AnyOf $s.OneOf -}}interface {
        {{$.Name}}TypeCheck()
    }

    type {{$.Name}}Intf interface {
        {{$.Name}}TypeCheck()
    }
    {{- else if eq $s.Type "array" -}}[]
        {{- if not $s.Items -}}{{- fail (printf "items arrays must have items: %s" $.Name) -}}{{- end -}}
        {{- template "type_name" (named $ "Item" $s.Items) -}}

        {{- /* Array properties embedded */ -}}
        {{- template "type_embedded" (named $ "Item" $s.Items) -}}

    {{- else if eq $s.Type "object" -}}
        {{- if $s.AdditionalProperties -}}map[string]
            {{- if not $s.AdditionalProperties.SchemaRef -}}{{fail "additionalItems must not be the bool case"}}{{- end -}}
            {{- /* Map properties normal */ -}}
            {{- template "type_name" (named $ "Item" $s.AdditionalProperties.SchemaRef) }}

            {{- /* Map properties embedded */ -}}
            {{- template "type_embedded" (named $ "Item" $s.AdditionalProperties.SchemaRef) -}}

        {{- else if $s.Properties}}struct {
            {{- /* Process regular struct fields */ -}}
            {{- range $name := keysReflect $s.Properties | sortAlpha -}}
                {{$element := index $s.Properties $name}}
                {{- if and (not $element.Ref) $element.Schema.Description -}}
                    {{- range $c := split "\n" (trim $element.Schema.Description)}}
    // {{$c}}
                    {{- end -}}
                {{- end}}
    {{camelcase $name}} {{template "type_name" (named $ $name $element)}} `json:"{{$name}}{{if not ($s.IsRequired $name)}},omitempty{{end}}"`
            {{- end}}
}

            {{- /* Process embedded structs */ -}}
            {{- range $name := keysReflect $s.Properties | sortAlpha -}}
                {{$element := index $s.Properties $name}}
                {{- template "type_embedded" (named $ (camelcase $name) $element) -}}
            {{- end -}}
        {{- end}}
    {{- else}}
        {{- primitive $ $s -}}
    {{- end -}}
{{- end -}}

{{- /*
Used to output a ref name, or the type itself
*/ -}}
{{- if .Object.Ref -}}
    {{- refName .Object.Ref -}}
{{- else -}}
    {{- template "schema_noref" (named $ "" .Object.Schema) -}}
{{- end -}}
