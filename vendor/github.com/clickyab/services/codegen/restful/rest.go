package restful

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"path/filepath"
	"text/template"

	"sort"

	"github.com/clickyab/services/array"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/codegen/annotate"
	"github.com/clickyab/services/codegen/plugins"
	"github.com/goraz/humanize"
	"golang.org/x/tools/imports"
)

const (
	tmpl = `
// Code generated build with restful DO NOT EDIT.

package {{ .Package }}

{{ range $c := .Data }}
// {{ $c.Comment }}
// @Route {
// 		url = {{ $c.URL }}
//		method = {{ $c.Method }}{{ if ne $c.Payload ""}}
//		payload = {{ $c.Payload }}{{end}}{{ if $c.Protected }}
//		middleware = authz.Authenticate{{end}}{{ if ne $c.Resources ""}}
//		resource = {{ $c.Resources }}{{end}}
{{ range $p := $c.Params}}//		{{ $p.Key }} = {{ $p.Value }}
{{end}}//		200 = {{ $c.Result }}
//		400 = controller.ErrorResponseSimple{{ if $c.Protected }}
//		401 = controller.ErrorResponseSimple{{end}}{{ if ne $c.Resources ""}}
//		403 = controller.ErrorResponseSimple{{end}}
// }
func ({{ $c.Variable }} {{ $c.Receiver }}) {{ $c.FuncName }}(ctx context.Context, w http.ResponseWriter, r *http.Request) {
{{ if ne $c.Payload ""}}pl := {{ $c.Variable }}.MustGetPayload(ctx).(*{{ $c.Payload }}){{end}}
	res, err := {{ $c.Variable }}.{{ $c.OriginalFuncName }}(ctx, r{{ if ne $c.Payload ""}}, pl{{end}})
	if err != nil {
		framework.Write(w, err, http.StatusBadRequest)
		return
	}
	framework.Write(w, res, http.StatusOK)
}
{{end}}
`
)

var tpl = template.Must(template.New("rest").Parse(tmpl))

type param struct {
	Key, Value string
}

type params []param

func (p params) Len() int {
	return len(p)
}

func (p params) Less(i, j int) bool {
	return p[i].Key < p[j].Key
}

func (p params) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type base struct {
	URL              string
	Method           string
	Protected        bool
	Resources        string
	Payload          string
	Result           string
	Variable         string
	Receiver         string
	OriginalFuncName string
	FuncName         string
	Comment          string
	Params           params
}

type context struct {
	Data    []base
	Package string
	File    string
}

type plugin struct {
}

func isContext(f humanize.File, ty humanize.Type) bool {
	alias := ""
	for i := range f.Imports {
		if f.Imports[i].Path == "context" {
			alias = f.Imports[i].Name
			if alias == "" {
				alias = "context"
			}
		}
	}

	if t, ok := ty.(*humanize.SelectorType); ok {
		if t.GetDefinition() == alias+".Context" {
			return true
		}
	}

	return false
}

func isRequest(f humanize.File, ty humanize.Type) bool {
	alias := ""
	for i := range f.Imports {
		if f.Imports[i].Path == "net/http" {
			alias = f.Imports[i].Name
			if alias == "" {
				alias = "http"
			}
		}
	}
	t, ok := ty.(*humanize.StarType)
	if !ok {
		return false
	}

	if t, ok := t.Target.(*humanize.SelectorType); ok {
		if t.GetDefinition() == alias+".Request" {
			return true
		}
	}

	return false
}

func isError(f humanize.File, ty humanize.Type) bool {
	// TODO : check for Error() string method
	if ty.GetDefinition() == "error" {
		return true
	}

	return false
}

func (p *plugin) GetOrder() int {
	return 99
}

func (p *plugin) GetType() []string {
	return []string{"Rest"}
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

func (p *plugin) Finalize(c interface{}, pkg *humanize.Package) error {
	ctx, ok := c.(context)
	if !ok {
		return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
	}

	out := &bytes.Buffer{}
	if err := tpl.Execute(out, ctx); err != nil {
		return err
	}
	res, err := imports.Process(ctx.File, out.Bytes(), nil)
	if err != nil {
		fmt.Println(out.String())
		return err
	}
	err = ioutil.WriteFile(ctx.File, res, 0644)
	if err != nil {
		return err
	}

	return appendToPkg(pkg, ctx.File)
}

func (p *plugin) FunctionIsSupported(f humanize.File, fn humanize.Function) bool {
	if fn.Receiver == nil {
		return false
	}
	if len(fn.Type.Results) != 2 {
		return false
	}
	if x := len(fn.Type.Parameters); x != 2 && x != 3 {
		return false
	}

	if !isContext(f, fn.Type.Parameters[0].Type) {
		return false
	}

	if !isRequest(f, fn.Type.Parameters[1].Type) {
		return false
	}

	if !isError(f, fn.Type.Results[1].Type) {
		return false
	}

	return true
}

func (p *plugin) ProcessFunction(c interface{}, pkg humanize.Package, f humanize.File, fn humanize.Function, an annotate.Annotate) (interface{}, error) {
	ctx, ok := c.(context)
	if !ok {
		ctx = context{}
	}
	rec := fn.Receiver.Type.GetDefinition()
	protected, _ := strconv.ParseBool(an.Items["protected"])
	res := an.Items["resource"]
	if res != "" {
		protected = true
	}

	var cn []string
	for i := range fn.Docs {
		d := strings.TrimLeft(fn.Docs[i], "/ ")
		if d[0] == '@' {
			break
		}
		cn = append(cn, d)
	}
	result := fn.Type.Results[0].Type.GetDefinition()
	if t, ok := fn.Type.Results[0].Type.(*humanize.StarType); ok {
		result = t.Target.GetDefinition()
	}
	meth := an.Items["method"]
	methUC := strings.ToUpper(meth[:1]) + strings.ToLower(meth[1:])
	name := strings.Split(fn.Name, ".")
	assert.True(len(name) == 2)
	data := base{
		URL:              an.Items["url"],
		OriginalFuncName: name[1],
		Method:           meth,
		FuncName:         name[1] + methUC,
		Receiver:         rec,
		Protected:        protected,
		Resources:        res,
		Variable:         fn.Receiver.Name,
		Result:           result,
		Comment:          strings.Join(cn, "\n"),
		Params:           make(params, 0),
	}

	for i := range an.Items {
		if !array.StringInArray(i, "method", "url", "protected", "resource") {
			data.Params = append(data.Params, param{Key: i, Value: an.Items[i]})
		}
	}
	sort.Sort(data.Params)

	if len(fn.Type.Parameters) == 3 {
		tn, ok := fn.Type.Parameters[2].Type.(*humanize.StarType)
		if !ok {
			return nil, fmt.Errorf("func %s 3rd parammeter is not a pointer, it should be", fn.Name)
		}
		data.Payload = tn.Target.GetDefinition()
	}

	ctx.Data = append(ctx.Data, data)
	ctx.Package = pkg.Name
	ctx.File = filepath.Join(filepath.Dir(f.FileName), pkg.Name+"_controllers.gen.go")
	return ctx, nil
}

func init() {
	plugins.Register(&plugin{})
}
