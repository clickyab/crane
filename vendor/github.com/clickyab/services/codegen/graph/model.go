package graph

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"time"

	"github.com/clickyab/services/codegen/annotate"
	"github.com/clickyab/services/codegen/plugins"
	"github.com/goraz/humanize"
	"golang.org/x/tools/imports"
)

type graph struct {
	pkg        humanize.Package
	file       humanize.File
	Ann        annotate.Annotate
	typ        humanize.TypeName
	View       *PermCode
	Epoch      string
	Type       string
	Key        string // number|time
	KeyType    string // time|int64
	URL        string
	Fill       string
	Entity     string
	Period     string // daily|weekly|monthly
	Scale      string // number|percent
	Layouts    []layout
	Conditions []condition
	format     []string
}

type condition struct {
	Field           string
	Filter          bool
	Searchable      bool
	Name            string
	Type            string
	Transform       string
	Title           string
	Perm            *PermCode
	HasPerm         bool
	FieldType       humanize.Type     `json:"-"`
	FieldTypeString string            `json:"-"`
	FilterValidMap  map[string]string `json:"filter_valid_map"`
	FilterValid     template.HTML     `json:"-"`
	Format          bool              `json:"-"`
	Data            string
}

type layout struct {
	Field     string
	FieldType string
	Hidden    bool
	Title     string
	Name      string
	Type      string
	Order     int64
	OmitEmpty bool
}

var errorNoGraphField = errors.New("not a graph field")

func handleFilter(p humanize.Package, f humanize.Field, mapPrefix string) (condition, error) {
	clm := condition{}
	tag := f.Tags.Get("json")
	if tag == "" {
		tag = f.Name
	}
	clm.Data = tag
	clm.Name = f.Name
	clm.Type = f.Tags.Get("type")
	clm.Searchable = strings.ToLower(f.Tags.Get("search")) == "true"
	clm.Filter = strings.ToLower(f.Tags.Get("filter")) == "true"

	clm.Transform = f.Tags.Get("map")
	if clm.Transform == "" && mapPrefix != "" {
		clm.Transform = mapPrefix + "." + tag
	}

	if clm.Filter && clm.Searchable {
		return clm, fmt.Errorf("both filter and search can not set on one field : %s", f.Name)
	}
	// Every thing is visible except when we note that
	clm.Title = f.Tags.Get("title")
	if clm.Title == "" {
		clm.Title = f.Name
	}
	if perm := f.Tags.Get("perm"); trim(perm) != "" {
		p, err := NewPermCode(perm)
		if err != nil {
			return clm, err
		}
		clm.Perm = p
		clm.HasPerm = true
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
func handleField(f humanize.Field) (layout, error) {
	l := layout{}
	if _, ok := f.Tags.Lookup("graph"); !ok {
		return l, errorNoGraphField
	}

	t := f.Tags.Get("graph")
	s := strings.Split(t, ",")
	if len(s) < 4 {
		return l, fmt.Errorf("graph tag is wrong for field %s", f.Name)
	}

	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	if i, ok := f.Type.(*humanize.IdentType); ok {
		l.FieldType = i.GetDefinition()
	} else {
		return l, fmt.Errorf("not supported type")
	}
	l.Field = f.Name
	l.Name = s[0]
	l.Title = s[1]
	l.Type = s[2]
	l.Hidden = strings.ToLower(s[3]) == "true"
	var err error
	l.Order, err = strconv.ParseInt(s[4], 10, 64)
	if err != nil {
		return l, err
	}
	if len(s) == 6 {
		l.OmitEmpty = s[5] == "true"
	}
	return l, nil
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

func (g graphPlugin) Finalize(c interface{}, p *humanize.Package) error {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	for i := range ctx {
		st := ctx[i].typ.Type.(*humanize.StructType)
		ls := make([]layout, 0)

		for _, f := range st.Fields {
			if !isExported(f.Name) {
				continue
			}
			if f.Name == ctx[i].Key {
				if x, ok := f.Type.(*humanize.SelectorType); ok {
					if x.GetDefinition() == "time.Time" {
						ctx[i].KeyType = "time.Time"
					}
				} else if x, ok := f.Type.(*humanize.IdentType); ok {
					switch x.GetDefinition() {
					case "int":
						ctx[i].KeyType = "int"
					case "int64":
						ctx[i].KeyType = "int64"
					default:
						return fmt.Errorf("type of key should be time.Time, int or int64")
					}
				} else {
					return fmt.Errorf("type of key should be time.Time, int or int64")
				}

			}

			l, err := handleField(*f)
			if err != nil {
				if err == errorNoGraphField {
					continue
				} else {
					return err
				}
			}

			ls = append(ls, l)
		}
		ctx[i].Layouts = ls
		ctx[i].Type = ctx[i].typ.Name

	}
	for i := range ctx {
		st := ctx[i].typ.Type.(*humanize.StructType)
		con := make([]condition, 0)
		mapPrefix := ctx[i].Ann.Items["map_prefix"]

		for _, f := range st.Fields {
			if !isExported(f.Name) {
				continue
			}
			c, err := handleFilter(*p, *f, mapPrefix)
			if err != nil {
				return err
			}
			c.Format = stringInArray(f.Name, ctx[i].format...)

			con = append(con, c)
		}
		ctx[i].Conditions = con
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
	f := filepath.Dir(p.Files[0].FileName)
	f = filepath.Join(f, "graphs.gen.go")
	res, err := imports.Process("", buf.Bytes(), nil)
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

	for i := range ctx {
		pp, err := humanize.ParsePackage(ctx[i].Ann.Items["controller"])

		if err != nil {
			return err
		}

		sorts := []string{}

		buf := &bytes.Buffer{}
		err = controller.Execute(buf, struct {
			Data                  graph
			PackageName           string
			ControllerPackageName string
			ValidSorts            template.HTML
			Model                 string
		}{
			Data:                  ctx[i],
			PackageName:           p.Name,
			ControllerPackageName: pp.Name,
			ValidSorts:            template.HTML(`"` + strings.Join(sorts, `","`) + `"`),
			Model:                 p.Path,
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
		//if err := appendToPkg(p, f); err != nil {
		//	return err
		//}

	}
	return nil
}

var (
	formater = regexp.MustCompile("Format([a-zA-Z]+)")
	prefix   = regexp.MustCompile("^_([a-zA-Z]+)")
	funcMap  = template.FuncMap{
		"getvar":   getVar,
		"scopeArg": scopeArg,
		"ucfirst":  ucFirst,
	}
	controller = template.Must(template.New("controller").Funcs(funcMap).Parse(ctrl))
	model      = template.Must(template.New("model").Funcs(funcMap).Parse(filterFunc))
)

func ucFirst(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func scopeArg(s string) template.HTML {
	switch s {
	case "parent":
		return template.HTML(`,permission.ScopeParent, permission.ScopeGlobal`)
	case "global":
		return template.HTML(`,permission.ScopeGlobal`)
	}
	return ""
}

type PermCode struct {
	Scope string
	Perm  string
	Total string
	Var   string
}

type context []graph

type graphPlugin struct {
}

func (g graphPlugin) GetOrder() int {
	return 10
}

func (g graphPlugin) GetType() []string {
	return []string{"Graph"}
}

func (g graphPlugin) StructureIsSupported(humanize.File, humanize.TypeName) bool {
	return true
}

func (g graphPlugin) ProcessStructure(
	c interface{},
	pkg humanize.Package,
	p humanize.File,
	f humanize.TypeName,
	a annotate.Annotate) (interface{},
	error) {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return nil, fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}
	period := "daily"
	if c := a.Items["period"]; c != "" {
		switch c {
		case "hourly":
			period = c
		case "daily":
			period = c
		case "weekly":
			period = c
		case "monthly":
			period = c
		default:
			return nil, errors.New("period should be one of hourly, daily, weekly or monthly")
		}
	}
	scale := "number"
	if c := a.Items["scale"]; c != "" {
		switch c {
		case "number":
			scale = c
		case "percent":
			scale = c
		default:
			return nil, errors.New("scale should be one of number or percent")
		}
	}
	epoch := "2018010100"
	if x := a.Items["epoch"]; x != "" {
		y, err := time.Parse("2006-01-02", x)
		if err != nil {
			return nil, err
		}
		epoch = y.Format("2006010215")
	}
	dt := graph{
		pkg:    pkg,
		file:   p,
		Ann:    a,
		typ:    f,
		Fill:   a.Items["fill"],
		Key:    a.Items["key"],
		Entity: a.Items["entity"],
		URL:    a.Items["url"],
		Epoch:  epoch,
		Period: period,
		Scale:  scale,
	}

	var err error
	dt.View, err = NewPermCode(a.Items["view"])
	if err != nil {
		return nil, err
	}

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

func isExported(s string) bool {
	if len(s) == 0 {
		panic("empty?")
	}

	return strings.ToUpper(s[:1]) == s[:1]
}

func trim(s string) string {
	return strings.Trim(s, " \n\t\"")
}
func init() {
	plugins.Register(&graphPlugin{})
}
