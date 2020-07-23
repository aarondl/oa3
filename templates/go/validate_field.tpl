{{- /* Top level object is a Schema, Name is the name of the local var */ -}}
ers = nil
{{- if $.Object.MaxLength}}
if err := support.MaxLength({{$.Name}}, {{$.Object.MaxLength}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MinLength}}
if err := support.MinLength({{$.Name}}, {{$.Object.MinLength}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.Maximum}}
if err := support.ValidateMax{{if eq $.Object.Type "integer"}}Int(int64({{else}}Float64(float64({{end}}{{$.Name}}), {{$.Object.Maximum}}, {{$.Object.ExclusiveMaximum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.Minimum}}
if err := support.ValidateMin{{if eq $.Object.Type "integer"}}Int(int64({{else}}Float64(float64({{end}}{{$.Name}}), {{$.Object.Minimum}}, {{$.Object.ExclusiveMinimum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MultipleOf}}
if err := support.ValidateMultipleOf{{if eq $.Object.Type "integer"}}Int(int64({{else}}Float64(float64({{end}}{{$.Name}}), {{$.Object.MultipleOf}}, {{$.Object.ExclusiveMinimum}}); err != nil {
    ers = append(ers, err)
}
{{- end -}}
{{- if $.Object.MaxItems}}
if err := support.MaxItems({{$.Name}}, {{$.Object.MaxItems}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MinItems}}
if err := support.MinItems({{$.Name}}, {{$.Object.MinItems}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MaxProperties}}
if err := support.MaxProperties({{$.Name}}, {{$.Object.MaxProperties}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if $.Object.MinProperties}}
if err := support.MinProperties({{$.Name}}, {{$.Object.MinProperties}}); err != nil {
    ers = append(ers, err)
}
{{end -}}
{{- if and $.Object.Enum (gt (len $.Object.Enum) 0) }}
if err := support.Enum({{$.Name}}, []string{
    {{- range $.Object.Enum -}}
        {{- if ne "string" (typeOf .) -}}
        {{- else -}}
        {{printf "%q" .}}, {{end -}}
    {{- end -}}
        }); err != nil {
    ers = append(ers, err)
}
{{end -}}
