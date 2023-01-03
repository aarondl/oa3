{{- $.Import "github.com/aarondl/oa3/openapi3spec" -}}
{{- $.Import "github.com/aarondl/oa3/support" -}}
var spec = {{codeForValue $.Spec}}
