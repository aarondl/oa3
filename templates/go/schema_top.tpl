{{- if $.Object.Description -}}
{{- range $c := split "\n" (trim $.Object.Description) -}}
// {{$c}}
{{end -}}
{{- end -}}
type {{title $.Name}} {{template "schema" .}}

{{- if and (not $.Object.Ref) (not $.Object.AnyOf) (not $.Object.OneOf) -}}
    {{- template "validate_schema" (recurseData $ "" $.Object) -}}
{{- end -}}
