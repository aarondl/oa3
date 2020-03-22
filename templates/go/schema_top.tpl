{{- if .Object.Description -}}
{{- range $c := split "\n" (trim .Object.Description) -}}
// {{$c}}
{{end -}}
{{- end -}}
type {{.Name}} {{template "schema" .}}
