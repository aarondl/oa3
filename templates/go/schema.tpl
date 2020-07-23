{{- /* Used to output a type name, takes a TemplateData with the object set to a schema ref */ -}}
{{- define "type_name" -}}
    {{- if .Object.Enum -}}
        {{- .Name -}}
    {{- else if or .Object.Ref (isInlinePrimitive .Object.Schema) -}}
        {{- template "schema" (recurseData $ .Name .Object) -}}
    {{- else -}}
        {{- if .Object.Schema.Nullable}}*{{end}}{{.Name}}
    {{- end -}}
{{- end -}}

{{- /* Outputs enum constants */ -}}
{{- define "type_enum" -}}
const (
    {{- range $value := $.Object.Enum}}
    {{$.Name}}{{camelcase $value}} {{$.Name}} = {{printf "%q" $value}}
    {{- end}}
)
{{- end -}}

{{- /* Used to output an embedded type, takes a TemplateData with the object set
to a schema ref */ -}}
{{- define "type_embedded" -}}
    {{- if $.Object.Enum}}
        {{template "type_enum" $}}
    {{- else if and (not .Object.Ref) (not (isInlinePrimitive .Object.Schema))}}

{{template "schema_top" $ -}}
    {{- end -}}
{{- end -}}

{{- /* Write out the schema after ensuring it's not a ref */ -}}
{{- define "schema_noref" -}}
    {{- $s := .Object -}}

    {{- if $s.Enum -}}string

    {{template "type_enum" $}}

    {{- else if or $s.AnyOf $s.OneOf -}}interface {
        {{$.Name}}TypeCheck()
    }

    type {{$.Name}}Intf interface {
        {{$.Name}}TypeCheck()
    }
    {{- else if eq $s.Type "array" -}}[]
        {{- if not $s.Items -}}{{- fail (printf "items arrays must have items: %s" $.Name) -}}{{- end -}}
        {{- template "type_name" (recurseData $ "Item" $s.Items) -}}

        {{- /* Array properties embedded */ -}}
        {{- template "type_embedded" (recurseData $ "Item" $s.Items) -}}

    {{- else if eq $s.Type "object" -}}
        {{- if $s.AdditionalProperties -}}map[string]
            {{- if not $s.AdditionalProperties.SchemaRef -}}{{fail "additionalItems must not be the bool case"}}{{- end -}}
            {{- /* Map properties normal */ -}}
            {{- template "type_name" (recurseData $ "Item" $s.AdditionalProperties.SchemaRef) }}

            {{- /* Map properties embedded */ -}}
            {{- template "type_embedded" (recurseData $ "Item" $s.AdditionalProperties.SchemaRef) -}}

        {{- else if $s.Properties}}struct {
            {{- /* Process regular struct fields */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and (not $element.Ref) $element.Schema.Description -}}
                    {{- range $c := split "\n" (trim $element.Schema.Description)}}
    // {{$c}}
                    {{- end -}}
                {{- end}}
    {{camelcase $name}} {{template "type_name" (recurseData $ $name $element)}} `json:"{{$name}}{{if not ($s.IsRequired $name)}},omitempty{{end}}"`
            {{- end}}
}

            {{- /* Process embedded structs */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- template "type_embedded" (recurseData $ (camelcase $name) $element) -}}
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
    {{- template "schema_noref" (recurseData $ "" .Object.Schema) -}}
{{- end -}}
