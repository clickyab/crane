package graph

const (
	filterFunc = `// Code generated build with graph DO NOT EDIT.

	package {{ .PackageName }}

type (
{{ range $m := .Data }}
	{{ $m.Type }}Array []{{ $m.Type }}
{{ end }}
)

{{ range $m := .Data }}

func ({{ $m.Type|getvar }}a {{ $m.Type }}Array) Filter(u permission.Interface){{ $m.Type }}Array {
	res := make({{ $m.Type }}Array, len({{ $m.Type|getvar }}a))
	for i := range {{ $m.Type|getvar }}a {
		res[i] = {{ $m.Type|getvar }}a[i].Filter(u)
	}

	return res
}

// Filter is for filtering base on permission
func ({{ $m.Type|getvar }} {{ $m.Type }}) Filter(u permission.Interface) {{ $m.Type }} {
	res := {{ $m.Type }}{}
	{{ range $clm := $m.Conditions }}
	{{ if not $clm.HasPerm }}res.{{ $clm.Name }} = {{ if $clm.Format }} {{ $m.Type|getvar }}.Format{{ $clm.Name}}(){{ else }}{{ $m.Type|getvar }}.{{ $clm.Name}}{{ end }}{{ end }}
	{{ end }}
	{{ range $clm := $m.Conditions }}
	{{ if $clm.HasPerm }}
	if _, ok := u.HasPermOn("{{ $clm.Perm.Perm }}", {{ $m.Type|getvar }}.OwnerID, {{ $m.Type|getvar }}.ParentID.Int64 {{ $clm.Perm.Scope|scopeArg }}); ok {
		res.{{ $clm.Name }} = {{ if $clm.Format }} {{ $m.Type|getvar }}.Format{{ $clm.Name}}()  {{ else }}{{ $m.Type|getvar }}.{{ $clm.Name}} {{ end }}
	}
	{{ end }}
	{{ end }}
	return res
}


func init () {
	{{ range $c:= $m.Conditions }}
		{{ if $c.Perm }}
		permission.RegisterPermission("{{ $c.Perm.Perm }}", "{{ $c.Perm.Perm }}");
		{{ end }}
	{{ end }}
}

{{end}}
`
)
