// URLBuilder builds a base url. Implementations are likely simple fixed
// strings or slightly more complicated variable replacement strings with
// defaults.
//
// Implementors:
{{- range $server := $.Spec.Servers}}
// - {{filterNonIdentChars $server.URL | title}}
{{- end}}
type URLBuilder interface {
	ToURL() string
}

// URL is a simple base url builder that's just a static string
type URL string
func (b URL) ToURL() string { return string(b) }

{{/* Called with an oa3 Server type */}}
{{define "urls_for_server"}}
{{- $server := $.Object -}}
{{- $serverName := filterNonIdentChars $server.URL | title -}}
{{- if $server.Description}}
// {{wrapWith 70 "\n// " (trimSuffix "\n" $server.Description)}}
{{- end -}}
	{{- if $server.Variables}}
type {{$serverName}} struct {
		{{- range $k := (keysReflect $server.Variables | sortAlpha)}}
			{{$k | title}} string
		{{- end}}
}
		{{- with $.Name}}
func ({{$serverName}}) {{$.Name}}ToURL() {}
		{{- end}}
func (s {{$serverName}}) ToURL() string {
		uri := `{{$server.URL}}`
		{{- range $k := (keysReflect $server.Variables | sortAlpha) -}}
			{{- $varname := $k | title -}}
			{{- $v := index $server.Variables $k -}}
			{{- $.Import "strings"}}
		if len(s.{{$varname}}) != 0 {
			{{- with $enumVals := $v.Enum}}
			switch s.{{$varname}} {
			case {{range $i, $enum := $enumVals}}{{if ne $i 0}}, {{end}}`{{$enum}}`{{end}}:
			default:
				panic("unknown server variable enum value: " + s.{{$varname}})
			}
			{{- end}}
			uri = strings.ReplaceAll(uri, `{{"{"}}{{$k}}{{"}"}}`, s.{{$varname}})
		} else {
			uri = strings.ReplaceAll(uri, `{{"{"}}{{$k}}{{"}"}}`, `{{$v.Default}}`)
		}
		{{- end}}
		return uri
}
	{{- else}}
var {{$serverName}} = URL{{$.Name}}(`{{$server.URL}}`)
	{{- end -}}
{{- end -}}

{{- range $server := $.Spec.Servers -}}
{{template "urls_for_server" (newData $ "" $server)}}
{{- end}}

{{- range $url, $path := $.Spec.Paths -}}
	{{- with $servers := $path.Servers -}}
		{{- if hasComplexServers $servers -}}
			{{- $intfName := $url | filterNonIdentChars | title}}
// URL{{$intfName}} is a simple string url like URLSimple
type URL{{$intfName}} string
func (b URL{{$intfName}}) ToURL() string { return string(b) }
func (b URL{{$intfName}}) {{$intfName}}Satisfy() {}

// URLBuilder{{$intfName}} builds a base url like URLBuilder
// but restricts the implementing types to a smaller subset.
//
// Implementors:
			{{- range $server := $servers}}
// - {{filterNonIdentChars $server.URL | title -}}
			{{- end}}
type URLBuilder{{$intfName}} interface {
	URLBuilder
	{{$intfName}}Satisfy()
}
			{{- range $server := $servers -}}
			{{template "urls_for_server" (newData $ $intfName $server)}}
			{{- end}}
		{{- end -}}
	{{- end -}}
    {{- range $method, $op := $path.Operations -}}
		{{- with $servers := $op.Servers -}}
			{{- if hasComplexServers $servers -}}
				{{- $intfName := printf "%s%s" ($url | filterNonIdentChars | title) ($method | filterNonIdentChars | title) }}
// URL{{$intfName}} is a simple url
type URL{{$intfName}} string
func (b URL{{$intfName}}) ToURL() string { return string(b) }
func (b URL{{$intfName}}) {{$intfName}}Satisfy() {}

// URLBuilder{{$intfName}} builds a base url like URLBuilder
// but restricts the implementing types to a smaller subset.
//
// Implementors:
				{{- range $server := $servers}}
// - {{filterNonIdentChars $server.URL | title -}}
				{{- end}}
type URLBuilder{{$intfName}} interface {
	URLBuilder
	{{$intfName}}Satisfy()
}
				{{- with $servers := $op.Servers -}}
					{{- range $server := $servers -}}
				{{template "urls_for_server" (newData $ $intfName $server)}}
					{{- end -}}
				{{- end -}}
			{{- end -}}
		{{- end -}}
	{{- end -}}
{{- end}}
