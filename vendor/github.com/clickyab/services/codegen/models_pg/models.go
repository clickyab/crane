package models

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/clickyab/services/codegen/annotate"
	"github.com/clickyab/services/codegen/plugins"

	"github.com/goraz/humanize"
	"github.com/jinzhu/inflection"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/imports"
)

type fieldModel struct {
	DB   string
	Name string
	Type string
}

type manyToMany struct {
	Base   string
	St1    string
	St2    string
	Field1 fieldModel
	Field2 fieldModel
}

type belongTo struct {
	Base   string
	St     string
	Field  fieldModel
	Target string
}

type hasMany struct {
	Base  string
	St    string
	Field fieldModel
}

type dataModel struct {
	Schema      string
	Table       string
	StructName  string
	FileName    string
	PackageName string
	FindBy      []fieldModel
	FilterBy    []fieldModel
	AutoIncr    bool
	Primaries   []string
	List        bool
	M2M         *manyToMany
	B2          *belongTo
	HM          *hasMany
	CreatedAt   *fieldModel
	UpdatedAt   *fieldModel
	Transaction []string
	Fields      string
}

type dataModels []dataModel

type context struct {
	data        map[string]dataModels
	p           humanize.Package
	packageName string
}
type modelsPlugin struct {
}

const tmpl = `
// Code generated build with models DO NOT EDIT.

package {{ .PackageName }}

{{ range $m := .Data }}
{{ if $m.Primaries }}
// Create{{ $m.StructName }} try to save a new {{ $m.StructName }} in database
func (m *Manager) Create{{ $m.StructName }}({{ $m.StructName|getvar }} ...*{{ $m.StructName }}) error {
	dum := make([]interface{}, len({{ $m.StructName|getvar }}))
	{{ if $m.CreatedAt }}now := time.Now(){{ else if $m.UpdatedAt }}now := time.Now(){{ end }}
	for i := range {{ $m.StructName|getvar }} {
		{{ if $m.CreatedAt }}{{ $m.StructName|getvar }}[i].CreatedAt = now{{ end }}
		{{ if $m.UpdatedAt }}{{ $m.StructName|getvar }}[i].UpdatedAt = now{{ end }}
		dum[i] =  {{ $m.StructName|getvar }}[i]
		if ii, ok := dum[i].(initializer.Simple); ok {
			ii.Initialize()
		}
    }

	return m.GetDbMap().Insert(dum...)
}

// Update{{ $m.StructName }} try to update {{ $m.StructName }} in database
func (m *Manager) Update{{ $m.StructName }}({{ $m.StructName|getvar }} ...*{{ $m.StructName }}) error {
	{{ if $m.UpdatedAt }}now := time.Now(){{ end }}
	dum := make([]interface{}, len({{ $m.StructName|getvar }}))
	for i := range {{ $m.StructName|getvar }} {
		{{ if $m.UpdatedAt }}{{ $m.StructName|getvar }}[i].UpdatedAt = now{{ end }}
		dum[i] =  {{ $m.StructName|getvar }}[i]
		if ii, ok := dum[i].(initializer.Simple); ok {
			ii.Initialize()
		}
    }

	_, err := m.GetDbMap().Update(dum...)
	return err
}
{{ end }}
{{ if $m.List }}

// List{{ $m.StructName|plural }}WithFilter try to list all {{ $m.StructName|plural }} without pagination
func (m *Manager) List{{ $m.StructName|plural }}WithFilter(filter string, params ...interface{}) []{{ $m.StructName }} {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []{{ $m.StructName }}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %s FROM %s %s",getSelectFields( {{ $m.StructName }}TableFull,""), {{ $m.StructName }}TableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}


// List{{ $m.StructName|plural }} try to list all {{ $m.StructName|plural }} without pagination
func (m *Manager) List{{ $m.StructName|plural }}() []{{ $m.StructName }} {
	return m.List{{ $m.StructName|plural }}WithFilter("")
}

// Count{{ $m.StructName|plural }}WithFilter count entity in {{ $m.StructName|plural }} table with valid where filter
func (m *Manager) Count{{ $m.StructName|plural }}WithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", {{ $m.StructName }}TableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// Count{{ $m.StructName|plural }} count entity in {{ $m.StructName|plural }} table
func (m *Manager) Count{{ $m.StructName|plural }}() int64 {
	return m.Count{{ $m.StructName|plural }}WithFilter("")
}

// List{{ $m.StructName|plural }}WithPaginationFilter try to list all {{ $m.StructName|plural }} with pagination and filter
func (m *Manager) List{{ $m.StructName|plural }}WithPaginationFilter(
	offset, perPage int,
	filter string,
	params ...interface{}) []{{ $m.StructName }} {
	var res []{{ $m.StructName }}
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(params)+1, len(params) +2)
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %s FROM %s %s",getSelectFields( {{ $m.StructName }}TableFull,""), {{ $m.StructName }}TableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// List{{ $m.StructName|plural }}WithPagination try to list all {{ $m.StructName|plural }} with pagination
func (m *Manager) List{{ $m.StructName|plural }}WithPagination(offset, perPage int) []{{ $m.StructName }} {
	return m.List{{ $m.StructName|plural }}WithPaginationFilter(offset, perPage, "")
}
{{ end }}
{{ range $by := $m.FindBy }}
// Find{{ $m.StructName }}By{{ $by.Name }} return the {{ $m.StructName }} base on its {{ $by.DB }}
func (m* Manager) Find{{ $m.StructName }}By{{ $by.Name }}({{ $by.DB|getvar }} {{ $by.Type }}) (*{{ $m.StructName }}, error) {
	var res {{ $m.StructName }}
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT %s FROM %s WHERE {{ $by.DB }}=$1",getSelectFields( {{ $m.StructName }}TableFull,""), {{ $m.StructName }}TableFull),
		{{ $by.DB|getvar }},
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
{{ end }}
{{ range $by := $m.FilterBy }}
// Filter{{ $m.StructName|plural }}By{{ $by.Name }} return all {{ $m.StructName|plural }} base on its {{ $by.DB }}, panic on query error
func (m *Manager) Filter{{ $m.StructName|plural }}By{{ $by.Name }}({{ $by.DB|getvar }} {{ $by.Type }}) []{{ $m.StructName }} {
	var res []{{ $m.StructName }}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %s FROM %s WHERE {{ $by.DB }}=$1",getSelectFields( {{ $m.StructName }}TableFull,""), {{ $m.StructName }}TableFull),
		{{ $by.DB|getvar }},
	)
	assert.Nil(err)

	return res
}
{{ end }}
{{ with $m.M2M }}

// Get{{ .St1|base }}{{ .St2|plural }} return all {{ .St2|plural }} belong to {{ .St1 }} (many to many with {{ .Base }})
func (m *Manager) Get{{ .St1|base }}{{ .St2|plural }}({{ .St1|getvar }} *{{ .St1 }}) []{{ .St2 }} {
	var res []{{ .St2 }}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT {{ .St2|getvar }}.* FROM %s {{ .Base|getvar }} JOIN %s {{ .St2|getvar }} ON {{ .Base|getvar }}.{{ .Field2.DB }} = {{ .St2|getvar }}.id WHERE {{ .Base|getvar }}.{{ .Field1.DB }}=$1",
			{{ .Base  }}TableFull,
			{{ .St2 }}TableFull,
		),
		{{ .St1|getvar }}.ID,
	)

	assert.Nil(err )
	return res
}

// Get{{ .St2|base }}{{ .St1|plural }} return all {{ .St1|plural }} belong to {{ .St2 }} (many to many with {{ .Base }})
func (m *Manager) Get{{ .St2|base }}{{ .St1|plural }}({{ .St2|getvar }} *{{ .St2 }}) []{{ .St1 }} {
	var res []{{ .St1 }}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT {{ .St1|getvar }}.* FROM %s {{ .Base|getvar }} JOIN %s {{ .St1|getvar }} ON {{ .Base|getvar }}.{{ .Field1.DB }} = {{ .St1|getvar }}.id WHERE {{ .Base|getvar }}.{{ .Field2.DB }}=$1",
			{{ .Base  }}TableFull,
			{{ .St1 }}TableFull,
		),
		{{ .St2|getvar }}.ID,
	)

	assert.Nil(err )
	return res
}

// Count{{ .St1|base }}{{ .St2|plural }} return count {{ .St2|plural }} belong to {{ .St1 }} (many to many with {{ .Base }})
func (m *Manager) Count{{ .St1|base }}{{ .St2|plural }}({{ .St1|getvar }} *{{ .St1 }}) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE {{ .Field1.DB }}=$1",
			{{ .Base  }}TableFull,
		),
		{{ .St1|getvar }}.ID,
	)

	assert.Nil(err )
	return res
}

// Count{{ .St2|base }}{{ .St1|plural }} return all {{ .St1|plural }} belong to {{ .St2 }} (many to many with {{ .Base }})
func (m *Manager) Count{{ .St2|base }}{{ .St1|plural }}({{ .St2|getvar }} *{{ .St2 }}) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE {{ .Field2.DB }}=$1",
			{{ .Base  }}TableFull,
		),
		{{ .St2|getvar }}.ID,
	)

	assert.Nil(err )
	return res
}
{{ end }}

{{ with $m.B2 }}
// Get{{ .St|base }}{{ .Base|plural }} return all {{ .Base|plural }} belong to {{ .St|base }}
func (m *Manager) Get{{ .St|base }}{{ .Base|plural }}({{ .St|getvar }} *{{ .St }}) []{{ .Base }} {
	return m.List{{ .Base|plural }}WithFilter("{{ .Field.DB }}=$1", {{.St|getvar}}.{{ .Target }})
}

// List{{ .St|base }}{{ .Base|plural }}WithPagination return all {{ .Base|plural }} belong to {{ .St|base }}
func (m *Manager) List{{ .St|base }}{{ .Base|plural }}WithPagination(offset, perPage int, {{ .St|getvar }} *{{ .St }}) []{{ .Base }} {
	return m.List{{ .Base|plural }}WithPaginationFilter(offset, perPage,"{{ .Field.DB }}=$1", {{.St|getvar}}.{{ .Target }})
}

// Count{{ .St|base }}{{ .Base|plural }} return count {{ .Base|plural }} belong to {{ .St|base }}
func (m *Manager) Count{{ .St|base }}{{ .Base|plural }}({{ .St|getvar }} *{{ .St }}) int64 {
	return m.Count{{ .Base|plural }}WithFilter("{{ .Field.DB }}=$1", {{ .St|getvar }}.{{ .Target }})
}
{{ end }}

{{ with $m.HM }}
// Get{{ .Base|base }}{{ .St|plural }} return all {{ .St|plural }} belong to {{ .Base|base }} (has many)
func (m *Manager) Get{{ .Base|base }}{{ .St|plural }}({{ .Base|getvar }} *{{ .Base }}) []{{ .St }} {
	var res []{{ .St }}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT %s FROM %s WHERE {{ .Field.DB }}=$1",
			getSelectFields( {{ .St }}TableFull,""),
			{{ .St }}TableFull,
		),
		{{.Base|getvar}}.ID,
	)

	assert.Nil(err)
	return res
}
{{ end }}

{{ range $byT := .Transaction }}
// Pre{{ $byT }} is gorp hook to prevent {{ $byT }} without transaction
func ({{ $m.StructName|getvar }} *{{ $m.StructName }}) Pre{{ $byT }}(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return errors.New("{{ $byT }} {{ $m.StructName }} must be in transaction")
	}
	return nil
}

{{ end }}
{{ end }}

`

const base = `
// Code generated build with models DO NOT EDIT.

package {{ .PackageName }}

import (
	"gopkg.in/gorp.v2"
)
const ({{ range $m := .Data }}
	// {{ $m.StructName }}Schema is the {{ $m.StructName }} module schema
	{{ $m.StructName }}Schema = "{{ $m.Schema }}"
	// {{ $m.StructName }}Table is the {{ $m.StructName }} table name
	{{ $m.StructName }}Table = "{{ $m.Table }}"
	// {{ $m.StructName }}TableFull is the {{ $m.StructName }} table name with schema
	{{ $m.StructName }}TableFull = {{ $m.StructName }}Schema + "." + {{ $m.StructName }}Table
{{ end }})

func getSelectFields(tb string, alias string) string {
	if alias != "" {
		alias +="."
	}
	switch tb {
{{ range $m := .Data }}
	case {{ $m.StructName }}TableFull:
		return fmt.Sprintf(` + "`{{ $m.Fields }}`" + `, alias)
{{ end }}
	}
	return ""
}

// Manager is the model manager for {{ .PackageName }} package
type Manager struct {
	model.Manager
}

// New{{ .PackageName|ucfirst }}Manager create and return a manager for this module
func New{{ .PackageName|ucfirst }}Manager() *Manager {
	return &Manager{}
}

// New{{ .PackageName|ucfirst }}ManagerFromTransaction create and return a manager for this module from a transaction
func New{{ .PackageName|ucfirst }}ManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize {{ .PackageName }} package
func (m *Manager) Initialize() {
{{ range $m := .Data }}
	m.AddTableWithNameAndSchema(
				{{ $m.StructName }}{},
				{{ $m.StructName }}Schema,
				{{ $m.StructName }}Table,
			){{ if $m.Primaries }}.SetKeys(
				{{ $m.AutoIncr }},
				{{ range $p := $m.Primaries }} "{{ $p }}",
				 {{end}}
			){{ end }}
{{ end }}
}
func init() {
	postgres.Register(New{{ .PackageName|ucfirst }}Manager())
}


`

var (
	funcMap = template.FuncMap{
		"ucfirst":    ucFirst,
		"md5":        md5Sum,
		"strip_type": stripType,
		"singular":   singular,
		"plural":     plural,
		"getvar":     getVar,
		"base":       getBaseType,
	}
	tpl  = template.Must(template.New("model").Funcs(funcMap).Parse(tmpl))
	tpl2 = template.Must(template.New("base").Funcs(funcMap).Parse(base))
)

func (a dataModels) Len() int {
	return len(a)
}
func (a dataModels) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a dataModels) Less(i, j int) bool {
	in := strings.Compare(a[i].Schema, a[j].Schema)
	if in == 0 {
		in = strings.Compare(a[i].Table, a[j].Table)
	}
	return in < 0
}

func ucFirst(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func md5Sum(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func singular(s string) string {
	return inflection.Singular(getBaseType(s))
}

func plural(s string) string {
	return inflection.Plural(getBaseType(s))
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
	if res == "m" {
		return "mm"
	}
	return res
}

func getBaseType(s string) string {
	arr := strings.Split(s, ".")
	if len(arr) == 1 {
		return arr[0]
	}
	if len(arr) == 2 {
		return arr[1]
	}

	return s
}

func trim(s string) string {
	return strings.Trim(s, " \n\t\"")
}

func stripType(in string) string {
	res := strings.Split(in, ".")
	if len(res) == 1 {
		return res[0]
	}

	if len(res) == 2 {
		return res[1]
	}

	logrus.Panic("invalid name")
	return ""
}

func returnErr(key string) (interface{}, error) {
	return nil, fmt.Errorf("the key %s is not exists", key)
}

func getFields(s *humanize.StructType) []string {
	var res []string
	for i := range s.Embeds {
		if t, ok := s.Embeds[i].Type.(*humanize.StructType); ok {
			res = append(res, getFields(t)...)
		}
	}

	for i := range s.Fields {
		tag := s.Fields[i].Tags
		t := tag.Get("db_transform")
		if t == "" {
			tm := strings.Split(tag.Get("db"), ",")
			if len(tm) > 0 && trim(tm[0]) != "" && trim(tm[0]) != "-" {
				t = `%[1]s"` + trim(tm[0]) + `"`
			}
		}
		if t == "" {
			tm := s.Fields[i].Name
			if tm[0] >= 'A' && tm[0] <= 'Z' {
				t = `%[1]s"` + tm + `"`
			}
		}
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (r *modelsPlugin) GetType() []string {
	return []string{"Model"}
}

func fieldMapper(f string, st *humanize.StructType, pkg humanize.Package) ([]fieldModel, error) {
	var res []fieldModel
	all := strings.Split(f, ",")
	for _, s := range all {
		s = strings.Trim(s, "\n\t\" ")
		if s != "" {
			fm := fieldModel{
				DB: s,
			}
			var ok bool
			fm.Name, fm.Type, ok = findFieldInStructure(s, st, pkg)
			if !ok {
				return nil, fmt.Errorf("can not find field for %s ", s)
			}
			res = append(res, fm)
		}
	}

	return res, nil
}

func getManyToMany(s string, pkg humanize.Package, st humanize.TypeName) (*manyToMany, error) {
	b := strings.Split(s, ",")
	if len(b) != 2 {
		return nil, fmt.Errorf("exactly two part needed, %s has %d part", s, len(b))
	}
	f1 := strings.Split(b[0], ":")
	f2 := strings.Split(b[1], ":")

	if len(f1) != 2 || len(f2) != 2 {
		return nil, fmt.Errorf("correct format is Table1:field_1, Table2: field_2")
	}

	st1 := trim(f1[0])
	mapper1, err := fieldMapper(trim(f1[1]), st.Type.(*humanize.StructType), pkg)
	if err != nil {
		return nil, err
	}
	if len(mapper1) != 1 {
		return nil, fmt.Errorf("field %s is not part of the structure", f1[1])
	}

	st2 := trim(f2[0])
	mapper2, err := fieldMapper(trim(f2[1]), st.Type.(*humanize.StructType), pkg)
	if err != nil {
		return nil, err
	}
	if len(mapper2) != 1 {
		return nil, fmt.Errorf("field %s is not part of the structure", f2[1])
	}

	return &manyToMany{
		Base:   st.Name,
		St1:    st1,
		St2:    st2,
		Field1: mapper1[0],
		Field2: mapper2[0],
	}, nil
}

func getBelongTo(s string, pkg humanize.Package, st humanize.TypeName) (*belongTo, error) {
	f := strings.Split(s, ":")
	if len(f) != 2 && len(f) != 3 {
		return nil, fmt.Errorf("correct format is Table:field_name")
	}

	b2 := trim(f[0])
	mapper, err := fieldMapper(trim(f[1]), st.Type.(*humanize.StructType), pkg)
	if err != nil {
		return nil, err
	}
	if len(mapper) != 1 {
		return nil, fmt.Errorf("field %s is not part of the structure", f[1])
	}
	var t string
	if len(f) == 2 {
		t = "ID"
	} else {
		t = f[2]
	}

	return &belongTo{
		Base:   st.Name,
		St:     b2,
		Field:  mapper[0],
		Target: t,
	}, nil
}

func getHasMany(s string, pkg humanize.Package, st humanize.TypeName) (*hasMany, error) {
	f := strings.Split(s, ":")
	if len(f) != 2 {
		return nil, fmt.Errorf("correct format is Table:field_name")
	}

	b2 := trim(f[0])
	mapper, err := fieldMapper(trim(f[1]), st.Type.(*humanize.StructType), pkg)
	if err != nil {
		return nil, err
	}
	if len(mapper) != 1 {
		return nil, fmt.Errorf("field %s is not part of the structure", f[1])
	}

	return &hasMany{
		Base:  st.Name,
		St:    b2,
		Field: mapper[0],
	}, nil
}

func findFieldInStructure(s string, st *humanize.StructType, pkg humanize.Package) (string, string, bool) {
	for _, fie := range st.Fields {
		if fie.Tags.Get("db") == s {
			// Bingo
			return fie.Name, fie.Type.GetDefinition(), true
		}
	}

	for _, emb := range st.Embeds {
		ident := emb.Type.(*humanize.IdentType)
		st, err := pkg.FindType(ident.Ident)
		if err == nil {
			if _, ok := st.Type.(*humanize.StructType); ok {
				res, tp, b := findFieldInStructure(s, st.Type.(*humanize.StructType), pkg)
				if b {
					return res, tp, b
				}
			}
		}
	}

	return "", "", false
}

func (r *modelsPlugin) ProcessStructure(
	c interface{},
	pkg humanize.Package,
	p humanize.File,
	f humanize.TypeName,
	a annotate.Annotate,
) (interface{}, error) {
	var ctx context
	var err error
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return nil, fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	} else {
		ctx.data = make(map[string]dataModels)
		ctx.p = pkg
		ctx.packageName = pkg.Name
	}

	data := dataModel{}
	var ok bool
	data.Schema, ok = a.Items["schema"]
	if !ok {
		return returnErr("@Model.schema")
	}

	data.Table, ok = a.Items["table"]
	if !ok {
		return returnErr("@Model.table")
	}

	list, ok := a.Items["list"]
	if ok && (list == "yes" || list == "true") {
		data.List = true
	}

	data.StructName = f.Name
	data.FileName = p.FileName

	data.PackageName = p.PackageName

	if fb, ok := a.Items["find_by"]; ok {
		data.FindBy, err = fieldMapper(fb, f.Type.(*humanize.StructType), pkg)
		if err != nil {
			return nil, err
		}
	}

	if fb, ok := a.Items["filter_by"]; ok {
		data.FilterBy, err = fieldMapper(fb, f.Type.(*humanize.StructType), pkg)
		if err != nil {
			return nil, err
		}
	}

	pr, ok := a.Items["primary"]
	if ok {
		pri := strings.Split(pr, ",")
		if len(pri) < 2 {
			return nil, fmt.Errorf("primary format is bool,key1 [,key2[,key3][,...]] got %s", pr)
		}
		data.AutoIncr, err = strconv.ParseBool(strings.Trim(pri[0], "\n\t\" "))
		if err != nil {
			return nil, err
		}

		for _, s := range pri[1:] {
			s = strings.Trim(s, "\n\t\" ")
			if s != "" {
				name, _, b := findFieldInStructure(s, f.Type.(*humanize.StructType), pkg)
				if !b {
					return nil, fmt.Errorf("can not find field for %s ", s)
				}
				data.Primaries = append(data.Primaries, name)
			}
		}
	}
	if m2m, ok := a.Items["many_to_many"]; ok {
		data.M2M, err = getManyToMany(m2m, pkg, f)
		if err != nil {
			return nil, err
		}
	}

	if b2, ok := a.Items["belong_to"]; ok {
		data.B2, err = getBelongTo(b2, pkg, f)
		if err != nil {
			return nil, err
		}
	}

	if hm, ok := a.Items["has_many"]; ok {
		data.HM, err = getHasMany(hm, pkg, f)
		if err != nil {
			return nil, err
		}
	}

	if byt, ok := a.Items["transaction"]; ok {
		for _, bt := range strings.Split(byt, ",") {
			if bt = trim(bt); bt != "" {
				data.Transaction = append(data.Transaction, ucFirst(bt))
			}
		}

	}

	if createdAt, err := fieldMapper("created_at", f.Type.(*humanize.StructType), pkg); err == nil {
		if len(createdAt) == 1 {
			data.CreatedAt = &createdAt[0]
		}
	}

	if updatedAt, err := fieldMapper("updated_at", f.Type.(*humanize.StructType), pkg); err == nil {
		if len(updatedAt) == 1 {
			data.UpdatedAt = &updatedAt[0]
		}
	}
	data.Fields = strings.Join(getFields(f.Type.(*humanize.StructType)), ",")
	ctx.data[p.FileName] = append(ctx.data[p.FileName], data)

	return ctx, nil
}
func appendToPkg(pkg *humanize.Package, f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	fl, err := humanize.ParseFile(string(b), pkg)
	if err != nil {
		return err
	}
	pkg.Files = append(pkg.Files, fl)
	return nil
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (r *modelsPlugin) Finalize(c interface{}, p *humanize.Package) error {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}
	//s, _ := json.MarshalIndent(ctx.data, "", "\t")
	//fmt.Print(string(s))

	var all dataModels
	last := ""
	for i := range ctx.data {
		sort.Sort(ctx.data[i])
		buf := &bytes.Buffer{}
		err := tpl.Execute(buf, struct {
			PackageName string
			Data        []dataModel
		}{
			PackageName: ctx.packageName,
			Data:        ctx.data[i],
		})
		if err != nil {
			return err
		}
		last = i
		f := i
		f = f[:len(f)-3] + ".gen.go"
		err = ioutil.WriteFile(f, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
		res, err := imports.Process(f, buf.Bytes(), nil)
		if err != nil {
			fmt.Println(buf.String())
			return err
		}

		err = ioutil.WriteFile(f, res, 0644)
		if err != nil {
			return err
		}
		if err := appendToPkg(p, f); err != nil {
			return err
		}

		all = append(all, ctx.data[i]...)
	}

	last = filepath.Join(filepath.Dir(last), ctx.packageName+".gen.go")
	buf := &bytes.Buffer{}
	sort.Sort(all)
	err := tpl2.Execute(buf, struct {
		PackageName string
		Data        dataModels
	}{
		PackageName: ctx.packageName,
		Data:        all,
	})
	if err != nil {
		return err
	}
	res, err := imports.Process(last, buf.Bytes(), nil)
	if err != nil {
		fmt.Println(buf.String())
		return err
	}
	err = ioutil.WriteFile(last, res, 0644)
	if err := appendToPkg(p, last); err != nil {
		return err
	}

	return err
}

func (r *modelsPlugin) StructureIsSupported(file humanize.File, fn humanize.TypeName) bool {
	return true
}

func (r *modelsPlugin) GetOrder() int {
	return 5
}

func init() {
	plugins.Register(&modelsPlugin{})
}
