package datatable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"services/codegen/annotate"
	"services/codegen/plugins"

	"github.com/goraz/humanize"
	"golang.org/x/tools/imports"
)

type dataTablePlugin struct {
}

type PermCode struct {
	Scope string
	Perm  string
	Total string
	Var   string
}

type dataTable struct {
	pkg     humanize.Package
	file    humanize.File
	Ann     annotate.Annotate
	typ     humanize.TypeName
	format  []string
	Type    string
	Column  []ColumnDef
	Columns template.HTML
	Actions map[string]*PermCode
	Fill    string
	View    *PermCode
	Entity  string
	URL     string
}

type context []dataTable

type ColumnDef struct {
	Data            string            `json:"data"`
	Name            string            `json:"name"`
	Searchable      bool              `json:"searchable"`
	Sortable        bool              `json:"sortable"`
	Visible         bool              `json:"visible"`
	Filter          bool              `json:"filter"`
	Title           string            `json:"title"`
	Format          bool              `json:"-"`
	Perm            *PermCode         `json:"-"`
	HasPerm         bool              `json:"-"`
	FieldType       humanize.Type     `json:"-"`
	FieldTypeString string            `json:"-"`
	FilterValidMap  map[string]string `json:"filter_valid_map"`
	FilterValid     template.HTML     `json:"-"`
	Transform       string            `json:"-"`
	Edit            *PermCode         `json:"-"`
}

var (
	formater = regexp.MustCompile("Format([a-zA-Z]+)")
	prefix   = regexp.MustCompile("^_([a-zA-Z]+)")
)

const (
	filterFunc = `
	// Code generated build with datatable DO NOT EDIT.

	package {{ .PackageName }}

type (
{{ range $m := .Data }}
	{{ $m.Type }}Array []{{ $m.Type }}
{{ end }}
)

{{ range $m := .Data }}

func ({{ $m.Type|getvar }}a {{ $m.Type }}Array) Filter(u base.PermInterface){{ $m.Type }}Array {
	res := make({{ $m.Type }}Array, len({{ $m.Type|getvar }}a))
	for i := range {{ $m.Type|getvar }}a {
		res[i] = {{ $m.Type|getvar }}a[i].Filter(u)
	}

	return res
}

// Filter is for filtering base on permission
func ({{ $m.Type|getvar }} {{ $m.Type }}) Filter(u base.PermInterface) {{ $m.Type }} {
	action :=[]string{}
	res := {{ $m.Type }}{}
	{{ range $clm := $m.Column }}
	{{ if not $clm.HasPerm }}res.{{ $clm.Name }} = {{ if $clm.Format }} {{ $m.Type|getvar }}.Format{{ $clm.Name}}(){{ else }}{{ $m.Type|getvar }}.{{ $clm.Name}}{{ end }}{{ end }}
	{{ end }}
	{{ range $clm := $m.Column }}
	{{ if $clm.Edit }}
	if _, ok := u.HasPermOn("{{ $clm.Edit.Perm }}", {{ $m.Type|getvar }}.OwnerID, {{ $m.Type|getvar }}.ParentID.Int64 {{ $clm.Edit.Scope|scopeArg }}); ok {
		action = append(action, "inline_{{$clm.Name}}")
	}
	{{ end }}
	{{ if $clm.HasPerm }}
	if _, ok := u.HasPermOn("{{ $clm.Perm.Perm }}", {{ $m.Type|getvar }}.OwnerID, {{ $m.Type|getvar }}.ParentID.Int64 {{ $clm.Perm.Scope|scopeArg }}); ok {
		res.{{ $clm.Name }} = {{ if $clm.Format }} {{ $m.Type|getvar }}.Format{{ $clm.Name}}()  {{ else }}{{ $m.Type|getvar }}.{{ $clm.Name}} {{ end }}
	}
	{{ end }}
	{{ end }}
	{{ range $act, $perm := $m.Actions }}
	if _, ok := u.HasPermOn("{{ $perm.Perm }}", {{ $m.Type|getvar }}.OwnerID, {{ $m.Type|getvar }}.ParentID.Int64 {{ $perm.Scope|scopeArg }}); ok {
		action = append(action, "{{ $act }}")
	}
	{{ end }}
	res.Actions = strings.Join(action, ",")
	return res
}


func init () {
	{{ range $act, $perm := $m.Actions }}
	base.RegisterPermission("{{ $perm.Perm }}", "{{ $perm.Perm }}");
	{{ end }}
	{{ range $c:= $m.Column }}
		{{ if $c.Perm }}
		base.RegisterPermission("{{ $c.Perm.Perm }}", "{{ $c.Perm.Perm }}");
		{{ end }}
		{{ if $c.Edit}}
		base.RegisterPermission("{{ $c.Edit.Perm }}", "{{ $c.Edit.Perm }}");
		{{ end }}
	{{ end }}
}

{{end}}
`

	controllerFunc = `
		// Code generated build with datatable DO NOT EDIT.


package {{ .ControllerPackageName }}


import (
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"gopkg.in/labstack/echo.v3"
	"{{ .Model  }}"
)

type list{{ .Data.Entity|ucfirst }}Response struct {
	Total   int64                  ` + "`json:\"total\"`" + `
	Data    {{ .PackageName }}.{{ .Data.Type }}Array ` + "`json:\"data\"`" + `
	Page    int                    ` + "`json:\"page\"`" + `
	PerPage int                    ` + "`json:\"per_page\"`" + `
	Definition base.Columns           ` + "`json:\"definition\"`" + `
}

var list{{ .Data.Entity|ucfirst }}Definition base.Columns

// @Route {
// 		url = {{ .Data.URL }}
//		method = get
//		_c_ = int , count per page
//		_p_ = int , page number
//		resource = {{ .Data.View.Total }}{{ if .HasSort }}
//		_sort_ = string, the sort and order like id:asc or id:desc available column {{ .ValidSorts }}{{end}}{{ range $f := .Data.Column }}{{ if $f.Filter }}
//		_{{ $f.Data }}_ = string , filter the {{ $f.Data }} field valid values are {{ $f.FilterValid }}{{ end }}{{ end }}{{ range $f := .Data.Column }}{{ if $f.Searchable }}
//		_{{ $f.Data }}_ = string , search the {{ $f.Data }} field {{ end }}{{ end }}
//		_def_ = bool, show definition in result?
//		200 = list{{ .Data.Entity|ucfirst }}Response
// }
func (u *Controller) list{{ .Data.Entity|ucfirst }}(ctx echo.Context) error {
	m :=  {{ .PackageName }}.New{{ .PackageName|ucfirst }}Manager()
	usr := authz.MustGetUser(ctx)
	p, c := httplib.GetPageAndCount(ctx.Request(), true)

	filter := make(map[string]string)
	{{ range $f := .Data.Column }}
	{{ if $f.Filter }}
	if e := ctx.Request().URL.Query().Get("{{ $f.Data }}"); e != "" && {{ $.PackageName }}.{{ $f.FieldTypeString }}(e).IsValid() {
		filter["{{ if ne $f.Transform "" }}{{ $f.Transform }}{{else}}{{ $f.Data }}{{end}}"] = e
	}
	{{ end }}
	{{ end }}
	search := make(map[string]string)
	{{ range $f := .Data.Column }}
	{{ if $f.Searchable }}
	if e := ctx.Request().URL.Query().Get("{{ $f.Data }}"); e != "" {
		search["{{ if ne $f.Transform "" }}{{ $f.Transform }}{{else}}{{ $f.Data }}{{end}}"] = e
	}
	{{ end }}
	{{ end }}
	{{ if .HasSort }}
	s := ctx.Request().URL.Query().Get("sort")
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		parts = append(parts, "asc")
	}
	sort := parts[0]
	if !array.StringInArray(sort, {{ .ValidSorts }}) {
		sort = ""
	}
	order := strings.ToUpper(parts[1])
	if !array.StringInArray(order, "ASC", "DESC") {
		order = "ASC"
	}
	{{ else }}
	sort := ""
	order := "ASC"
	{{ end }}

	params := make(map[string]string)
	for _, i := range ctx.ParamNames() {
		params[i] = ctx.Param(i)
	}

	pc := base.NewPermInterfaceComplete(usr, usr.ID, "{{ .Data.View.Perm }}", "{{ .Data.View.Scope }}")
	dt, cnt := m.{{ .Data.Fill }}(pc, filter, search, params, sort, order, p, c)
	res := 		list{{ .Data.Entity|ucfirst }}Response{
		Total:   cnt,
		Data:    dt.Filter(usr),
		Page:    p,
		PerPage: c,
	}
	if ctx.Request().URL.Query().Get("def") == "true" {
		res.Definition = list{{ .Data.Entity|ucfirst }}Definition
	}
	return u.OKResponse(
		ctx,
		res,
	)
}

func init() {
	tmp := []byte(` + "` {{ .Data.Columns }} `" + `)
	assert.Nil(json.Unmarshal(tmp, &list{{ .Data.Entity|ucfirst }}Definition))
}

`
)

var (
	funcMap = template.FuncMap{
		"getvar":   getVar,
		"scopeArg": scopeArg,
		"ucfirst":  ucFirst,
	}
	model      = template.Must(template.New("model").Funcs(funcMap).Parse(filterFunc))
	controller = template.Must(template.New("controller").Funcs(funcMap).Parse(controllerFunc))
)

func scopeArg(s string) template.HTML {
	switch s {
	case "parent":
		return template.HTML(`,base.ScopeParent, base.ScopeGlobal`)
	case "global":
		return template.HTML(`,base.ScopeGlobal`)
	}
	return ""
}

func ucFirst(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func trim(s string) string {
	return strings.Trim(s, " \n\t\"")
}

func getVar(s string) string {
	if str := strings.ToLower(trim(s)); len(str) < 3 {
		return str
	}
	s = CamelToSnake(s)
	arr := strings.Split(strings.ToLower(s), "_")
	res := ""
	for _, i := range arr {
		i = strings.Trim(i, " \n\t\"")
		if i != "" {
			res += i[0:1]
		}
	}

	return res
}

func NewPermCode(s string) (*PermCode, error) {
	parts := strings.Split(s, ":")
	if len(parts) == 1 {
		parts = append(parts, "global")
	}
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid perm string %s", s)
	}

	res := &PermCode{
		Scope: parts[1],
		Perm:  parts[0],
		Total: s,
		Var:   getVar(strings.Join(parts, "_")),
	}

	return res, nil
}

func isExported(s string) bool {
	if len(s) == 0 {
		panic("empty?")
	}

	return strings.ToUpper(s[:1]) == s[:1]
}

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (e dataTablePlugin) GetType() []string {
	return []string{"DataTable"}
}

func getDoc(d humanize.Docs, name string) string {
	var res []string
	for _, s := range d {
		s := strings.Trim(s, " /")
		if s[0] == '@' {
			break
		}
		res = append(res, s)
	}
	if len(res) > 1 {
		return res[1]
	}
	return name
}

func extractValidFilter(p humanize.Package, t humanize.Type) (map[string]string, string) {
	id, ok := t.(*humanize.IdentType)
	if !ok {
		return map[string]string{}, ""
	}
	res := make(map[string]string)
	comma := []string{}
	for i := range p.Files {
		for _, c := range p.Files[i].Constants {
			if c.Type.GetDefinition() == id.Ident {
				v := strings.Trim(c.Value, "\"`")
				res[v] = getDoc(c.Docs, c.Name)
				comma = append(comma, v)
			}
		}
	}

	return res, `"` + strings.Join(comma, `","`) + `"`
}

func handleField(p humanize.Package, f humanize.Field, mapPrefix string) (ColumnDef, error) {
	clm := ColumnDef{}
	tag := f.Tags.Get("json")
	if tag == "" {
		tag = f.Name
	}
	clm.Data = tag
	clm.Name = f.Name
	clm.Searchable = strings.ToLower(f.Tags.Get("search")) == "true"
	clm.Sortable = strings.ToLower(f.Tags.Get("sort")) == "true"
	clm.Filter = strings.ToLower(f.Tags.Get("filter")) == "true"

	clm.Transform = f.Tags.Get("map")
	if clm.Transform == "" && mapPrefix != "" {
		clm.Transform = mapPrefix + "." + tag
	}

	if clm.Filter && clm.Searchable {
		return ColumnDef{}, fmt.Errorf("both filter and search can not set on one field : %s", f.Name)
	}
	// Every thing is visible except when we note that
	clm.Visible = strings.ToLower(f.Tags.Get("visible")) != "false"
	clm.Title = f.Tags.Get("title")
	if clm.Title == "" {
		clm.Title = f.Name
	}
	if perm := f.Tags.Get("perm"); trim(perm) != "" {
		p, err := NewPermCode(perm)
		if err != nil {
			return ColumnDef{}, err
		}
		clm.Perm = p
		clm.HasPerm = true
	}

	if edit := f.Tags.Get("edit"); trim(edit) != "" {
		p, err := NewPermCode(edit)
		if err != nil {
			return ColumnDef{}, err
		}
		clm.Edit = p
	}

	clm.FieldType = f.Type
	if ii, ok := f.Type.(*humanize.IdentType); ok {
		clm.FieldTypeString = ii.Ident
	}
	if clm.Filter {
		var tmp string
		clm.FilterValidMap, tmp = extractValidFilter(p, f.Type)
		clm.FilterValid = template.HTML(tmp)
	}

	return clm, nil
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (e dataTablePlugin) Finalize(c interface{}, p humanize.Package) error {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	for i := range ctx {
		res := make(map[string]*PermCode)
		for key := range ctx[i].Ann.Items {
			if prefix.MatchString(strings.Trim(key, " ")) {
				var err error
				res[key[1:]], err = NewPermCode(ctx[i].Ann.Items[key])
				if err != nil {
					return err
				}
			}
		}
		ctx[i].Actions = res
		columns := make([]ColumnDef, 0)
		st := ctx[i].typ.Type.(*humanize.StructType)
		mapPrefix := ctx[i].Ann.Items["map_prefix"]
		for _, f := range st.Fields {
			if isExported(f.Name) && f.Tags.Get("json") != "-" {
				clm, err := handleField(p, *f, mapPrefix)
				if err != nil {
					return err
				}
				clm.Format = stringInArray(f.Name, ctx[i].format...)
				columns = append(columns, clm)
			}
		}

		for _, fe := range st.Embeds {
			tE, err := p.FindType(fe.Type.(*humanize.IdentType).Ident)
			if err != nil {
				return err
			}
			for _, f := range tE.Type.(*humanize.StructType).Fields {
				if isExported(f.Name) && f.Tags.Get("json") != "-" {
					clm, err := handleField(p, *f, mapPrefix)
					if err != nil {
						return err
					}
					clm.Format = stringInArray(f.Name, ctx[i].format...)
					columns = append(columns, clm)
				}
			}
		}
		ctx[i].Column = columns
		ctx[i].Type = ctx[i].typ.Name
		j, err := json.MarshalIndent(columns, "\t", "\t")
		if err != nil {
			return err
		}
		ctx[i].Columns = template.HTML(string(j))
	}
	buf := &bytes.Buffer{}
	err := model.Execute(buf, struct {
		Data        context
		PackageName string
	}{
		ctx,
		p.Name,
	})
	if err != nil {
		return err
	}
	model := p.Path
	f := filepath.Dir(p.Files[0].FileName)
	f = filepath.Join(f, "datatables.gen.go")
	res, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		fmt.Println(buf.String())
		return err
	}

	err = ioutil.WriteFile(f, res, 0644)
	if err != nil {
		return err
	}

	for i := range ctx {
		pp, err := humanize.ParsePackage(ctx[i].Ann.Items["controller"])
		if err != nil {
			return err
		}

		sorts := []string{}
		for _, j := range ctx[i].Column {
			if j.Sortable {
				sorts = append(sorts, j.Data)
			}
		}

		buf := &bytes.Buffer{}
		err = controller.Execute(buf, struct {
			Data                  dataTable
			PackageName           string
			ControllerPackageName string
			ValidSorts            template.HTML
			HasSort               bool
			Model                 string
		}{
			Data:                  ctx[i],
			PackageName:           p.Name,
			ControllerPackageName: pp.Name,
			ValidSorts:            template.HTML(`"` + strings.Join(sorts, `","`) + `"`),
			HasSort:               len(sorts) > 0,
			Model:                 model,
		})
		if err != nil {
			return err
		}
		f := filepath.Dir(pp.Files[0].FileName)
		f = filepath.Join(f, ctx[i].Ann.Items["entity"]+"_controller.gen.go")
		res, err := imports.Process("", buf.Bytes(), nil)
		if err != nil {
			fmt.Println(buf.String())
			return err
		}

		err = ioutil.WriteFile(f, res, 0644)
		if err != nil {
			return err
		}
	}

	//j, _ := json.MarshalIndent(ctx[0].Column, "\t", "\t")
	//fmt.Println(string(j))
	return nil
}

func (r *dataTablePlugin) ProcessStructure(
	c interface{},
	pkg humanize.Package,
	p humanize.File,
	f humanize.TypeName,
	a annotate.Annotate,
) (interface{}, error) {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return nil, fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	dt := dataTable{
		pkg:  pkg,
		file: p,
		Ann:  a,
		typ:  f,
		Fill: a.Items["fill"],
	}
	var err error
	dt.View, err = NewPermCode(a.Items["view"])
	if err != nil {
		return nil, err
	}

	dt.Entity = a.Items["entity"]
	dt.URL = a.Items["url"]

	for i := range pkg.Files {
		for _, fn := range pkg.Files[i].Functions {
			if fn.Receiver != nil {
				rec := fn.Receiver.Type
				if s, ok := rec.(*humanize.StarType); ok {
					rec = s.Target
				}
				if f.Name == rec.GetDefinition() {
					// found a function
					res := formater.FindStringSubmatch(fn.Name)
					if len(res) == 2 {
						dt.format = append(dt.format, res[1])
					}
				}
			}
		}
	}

	ctx = append(ctx, dt)
	return ctx, nil
}

func (r *dataTablePlugin) StructureIsSupported(file humanize.File, fn humanize.TypeName) bool {
	return true
}

func (r *dataTablePlugin) GetOrder() int {
	return 10
}

func init() {
	plugins.Register(&dataTablePlugin{})
}
