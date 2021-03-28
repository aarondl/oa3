// {{.Name}} is a client package to interact with the api.
{{- if $.Spec.Info.Description}}
//
// {{wrapWith 70 "\n// " (trimSuffix "\n" $.Spec.Info.Description)}}
{{- end}}
export default class {{.Name}} {
    baseUrl: string;

    constructor(baseUrl: string) {
        if (baseUrl === null) {
            this.baseUrl = '{{(index $.Spec.Servers 0).URL}}';
        } else {
            this.baseUrl = baseUrl;
        }
    }
{{range $url, $path := $.Spec.Paths -}}
{{- range $method, $op := $path.Operations -}}
{{- $opname := lowerFirst $op.OperationID }}
    // {{$opname}} {{$method}} {{$url}}
    {{if $op.Description -}}
    // {{wrapWith 70 "\n    // " (trimSuffix "\n" $op.Description)}}
    {{end -}}
    {{$opname}}(
        {{- if $op.RequestBody -}}body: any{{- end -}}
        {{- range $i, $param := $op.Parameters -}}
        {{- if or (ne $i 0) $op.RequestBody -}}, {{end -}}
        {{- lowerFirst $param.Name -}}: {{if eq $param.Schema.Type "integer"}}{{"number"}}{{else}}{{$param.Schema.Type}}{{end}}
        {{- end -}}
    ): Promise<Response> {
        let url = '{{$url}}';
        {{- range $i, $param := $op.Parameters -}}
            {{- if eq "path" $param.In}}
        url = url.replace('{{"{"}}{{$param.Name}}{{"}"}}', {{lowerFirst $param.Name -}}.toString());
            {{- end -}}
        {{- end -}}

        {{- $hasQuery := false -}}
        {{- range $i, $param := $op.Parameters -}}
            {{- if eq "query" $param.In -}}
                {{- if not $hasQuery -}}
                    {{- $hasQuery = true}}

        let query = '';
                {{- end}}
        if ({{lowerFirst $param.Name}} !== undefined) {
            if (query.length != 0) { query += '&' }
            query += '{{$param.Name}}=' + encodeURIComponent({{lowerFirst $param.Name}}.toString());
        }
                {{- if $param.Required -}}
        {{" "}}else {
            throw '{{lowerFirst $param.Name}} cannot be undefined'
        }
                {{- end -}}
            {{- end -}}
        {{- end -}}
        {{- if $hasQuery}}
        if (query.length != 0) {
            query = '?' + query;
        }
        {{- end}}

        let headers = new Headers();
        headers.set('Content-Type', 'application/json');
        {{- range $i, $param := $op.Parameters -}}
            {{- if eq "header" $param.In}}
        headers.set('{{$param.Name}}', {{lowerFirst $param.Name -}}.toString());
            {{- end -}}
        {{- end}}

        const params = {
            method: '{{upper $method}}',
            headers: headers,
            {{- if $op.RequestBody}}
            body: JSON.stringify(body),
            {{- end}}
        };

        return fetch(new Request(this.baseUrl + url{{if $hasQuery}} + query{{end}}, params));
    }
{{end -}}
{{- end -}}
}
