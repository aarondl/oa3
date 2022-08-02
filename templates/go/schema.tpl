{{- /* Used to output a type name, takes a TemplateData with the object set to a schema ref */ -}}
{{- define "type_name" -}}
    {{- if or $.Object.Ref (isInlinePrimitive .Object.Schema) -}}
        {{- if $.Object.Enum -}}
            {{- if $.Object.Ref -}}
                {{- refName $.Object.Ref -}}
            {{- else -}}
                {{title .Name}}
            {{- end -}}
        {{- else -}}
            {{- template "schema" (recurseData $ $.Name $.Object) -}}
        {{- end -}}
    {{- else -}}
        {{- if $.Object.Schema.Nullable}}*{{end}}{{.Name}}
    {{- end -}}
{{- end -}}

{{- /* Outputs enum constants, takes a Schema (not a ref) */ -}}
{{- define "type_enum" -}}
{{if or ($.Object.Nullable) (not $.Required)}}var{{else}}const{{end}} (
    {{- range $value := $.Object.Enum}}
    {{title $.Name}}{{filterNonIdentChars $value | snakeToCamel | title }} = {{title $.Name}}({{printf "%q" $value}})
    {{- end}}
)
{{- end -}}

{{- /* Used to output an embedded type, takes a TemplateData with the object set
to a schema ref */ -}}
{{- define "type_embedded" -}}
    {{- if and (not $.Object.Ref) $.Object.Enum}}

{{template "schema_top" (newDataRequired $ $.Name $.Object true)}}
    {{- else if and (not .Object.Ref) (not (isInlinePrimitive .Object.Schema))}}

{{template "schema_top" $ -}}
    {{- end -}}
{{- end -}}

{{- /* Write out the schema after ensuring it's not a ref */ -}}
{{- define "schema_noref" -}}
    {{- $s := .Object -}}

    {{- if $s.Enum -}}string

    {{template "type_enum" (newDataRequired $ $.Name $s true)}}

    {{- else if or $s.AnyOf $s.OneOf -}}interface {
        {{$.Name}}TypeCheck()
    }

    type {{$.Name}}Intf interface {
        {{$.Name}}TypeCheck()
    }
    {{- else if eq $s.Type "array" -}}[]
        {{- template "type_name" (recurseDataSetRequired $ "Item" $s.Items true) -}}

        {{- /* Array properties embedded */ -}}
        {{- template "type_embedded" (recurseDataSetRequired $ "Item" $s.Items true) -}}

    {{- else if eq $s.Type "object" -}}
        {{- if or (eq 0 (len $s.Properties)) $s.AdditionalProperties -}}map[string]
            {{- if $s.AdditionalProperties -}}
                {{- if not $s.AdditionalProperties.SchemaRef -}}{{fail "additionalItems must not be the bool case"}}{{- end -}}
                {{- /* Map properties normal */ -}}
                {{- template "type_name" (recurseDataSetRequired $ "Item" $s.AdditionalProperties.SchemaRef true) }}

                {{- /* Map properties embedded */ -}}
                {{- template "type_embedded" (recurseDataSetRequired $ "Item" $s.AdditionalProperties.SchemaRef true) -}}
            {{- else -}}
            interface{}
            {{- end -}}

        {{- else if $s.Properties}}struct {
            {{- /* Process regular struct fields */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- if and (not $element.Ref) $element.Schema.Description -}}
                    {{- range $c := split "\n" (trim $element.Schema.Description)}}
    // {{$c}}
                    {{- end -}}
                {{- end -}}
                {{- $elementRequired := $s.IsRequired $name -}}
                {{- $shouldWrap := or $element.Ref $element.Schema.Enum (and (ne $element.Type "array") (ne $element.Type "object"))}}
    {{camelcase $name}} {{if and $shouldWrap ($element.Schema.Nullable) (not $elementRequired) -}}
                    {{- $.Import "github.com/aarondl/opt/omitnull" -}}
                    omitnull.Val[{{template "type_name" (recurseDataSetRequired $ (camelcase $name) $element $elementRequired)}}]
                {{- else if and $shouldWrap ($element.Schema.Nullable) $elementRequired -}}
                    {{- $.Import "github.com/aarondl/opt/null" -}}
                    null.Val[{{template "type_name" (recurseDataSetRequired $ (camelcase $name) $element $elementRequired)}}]
                {{- else if and $shouldWrap (not $element.Schema.Nullable) (not $elementRequired) -}}
                    {{- $.Import "github.com/aarondl/opt/omit" -}}
                    omit.Val[{{template "type_name" (recurseDataSetRequired $ (camelcase $name) $element $elementRequired)}}]
                {{- else -}}
                    {{template "type_name" (recurseDataSetRequired $ (camelcase $name) $element $elementRequired)}}
                {{- end}} `json:"{{$name}}{{if not $elementRequired}},omitempty{{end}}"`
            {{- end}}
}

            {{- /* Process embedded structs */ -}}
            {{- range $name, $element := $s.Properties -}}
                {{- template "type_embedded" (recurseDataSetRequired $ (camelcase $name) $element ($s.IsRequired $name)) -}}
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
