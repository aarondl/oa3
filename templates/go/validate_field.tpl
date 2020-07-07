{{- /* Top level object is a Schema, Name is the name of the local var */ -}}
{{- $name := $.Name -}}
{{- if $.Object.Nullable -}}
{{- $name = printf "%s.%s" $.Name (primitive $ $.Object | replace "null." "") -}}
{{- end -}}

ers = nil
{{- if $.Object.MaxLength}}
if err := support.ValidateMaxLength(string({{$name}}), {{$.Object.MaxLength}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MinLength}}
if err := support.ValidateMinLength(string({{$name}}), {{$.Object.MinLength}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.Maximum}}
if err := support.ValidateMax{{if eq $.Object.Type "integer"}}Int(int64({{else}}Float64(float64({{end}}{{$name}}), {{$.Object.Maximum}}, {{$.Object.ExclusiveMaximum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.Minimum}}
if err := support.ValidateMin{{if eq $.Object.Type "integer"}}Int(int64({{else}}Float64(float64({{end}}{{$name}}), {{$.Object.Minimum}}, {{$.Object.ExclusiveMinimum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MultipleOf}}
if err := support.ValidateMultipleOf{{if eq $.Object.Type "integer"}}Int(int64({{else}}Float64(float64({{end}}{{$name}}), {{$.Object.MultipleOf}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MaxItems}}
if err := support.ValidateMaxItems({{$name}}, {{$.Object.MaxItems}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MinItems}}
if err := support.ValidateMinItems({{$name}}, {{$.Object.MinItems}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MaxProperties}}
if err := support.ValidateMaxProperties({{$name}}, {{$.Object.MaxProperties}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MinProperties}}
if err := support.ValidateMinProperties({{$name}}, {{$.Object.MinProperties}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.Pattern }}
if err := support.ValidatePattern(string({{$name}}), {{printf $.Object.Pattern}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if and $.Object.Enum (gt (len $.Object.Enum) 0) }}
if err := support.ValidateEnum(string({{$name}}), []string{
    {{- range $i, $v := $.Object.Enum -}}
        {{- if ne "string" (typeOf .) -}}
        {{- else -}}
        {{printf "%q" $v}}{{if gt $i 0}}, {{end}}{{end -}}
    {{- end -}}
        }); err != nil {
    ers = append(ers, err)
}
{{end -}}
