{{- /* Top level object is a Schema, Name is the name of the local var */ -}}
{{- $name := $.Name}}
ers = nil
{{- if $.Object.MaxLength}}
if err := support.ValidateMaxLength({{$name}}, {{$.Object.MaxLength}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MinLength}}
if err := support.ValidateMinLength({{$name}}, {{$.Object.MinLength}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if and $.Object.Maximum (or (eq $.Object.Type "number") (eq $.Object.Type "integer"))}}
if err := support.ValidateMaxNumber({{$name}}, {{$.Object.Maximum}}, {{$.Object.ExclusiveMaximum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if and $.Object.Maximum (eq $.Object.Type "string") (eq (printf $.Object.Format) "decimal")}}
if err := support.ValidateMaxShopspringDecimal({{$name}}, decimal.NewFromFloat({{$.Object.Maximum}}), {{$.Object.ExclusiveMaximum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if and $.Object.Minimum (or (eq $.Object.Type "number") (eq $.Object.Type "integer"))}}
if err := support.ValidateMinNumber({{$name}}, {{$.Object.Minimum}}, {{$.Object.ExclusiveMinimum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if and $.Object.Minimum (eq $.Object.Type "string") (eq (printf $.Object.Format) "decimal")}}
if err := support.ValidateMinShopspringDecimal({{$name}}, decimal.NewFromFloat({{$.Object.Minimum}}), {{$.Object.ExclusiveMinimum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MultipleOf}}
if err := support.ValidateMultipleOf{{if eq $.Object.Type "integer"}}Int{{else}}Float{{end}}({{$name}}, {{$.Object.MultipleOf}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MaxItems}}
if err := support.ValidateMaxItems({{$name}}, {{$.Object.MaxItems}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MinItems}}
if err := support.ValidateMinItems({{$name}}, {{$.Object.MinItems}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MaxProperties}}
if err := support.ValidateMaxProperties({{$name}}, {{$.Object.MaxProperties}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MinProperties}}
if err := support.ValidateMinProperties({{$name}}, {{$.Object.MinProperties}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.Pattern }}
if err := support.ValidatePattern({{$name}}, {{printf $.Object.Pattern}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if and $.Object.Format (eq (printf $.Object.Format) "uuid") }}
if err := support.ValidateFormatUUIDv4(string({{$name}})); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if and $.Object.Format (eq (printf $.Object.Format) "decimal") }}
if err := support.ValidateFormatDecimal({{$name}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if and $.Object.Enum (gt (len $.Object.Enum) 0) }}
if err := support.ValidateEnum({{$name}}, []string{
    {{- range $i, $v := $.Object.Enum -}}
        {{- if ne "string" (typeOf .) -}}
        {{- else -}}
        {{if gt $i 0}}, {{end}}{{printf "%q" $v}}{{end -}}
    {{- end -}}
        }); err != nil {
    ers = append(ers, err)
}
{{end -}}
