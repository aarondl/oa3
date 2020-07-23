{{- if $.Object.Description -}}
{{- range $c := split "\n" (trim $.Object.Description) -}}
// {{$c}}
{{end -}}
{{- end -}}
type {{$.Name}} {{template "schema" .}}

{{- if not $.Object.Ref -}}
    {{- template "validate_schema" (named $ "" $.Object) -}}
{{- end -}}
