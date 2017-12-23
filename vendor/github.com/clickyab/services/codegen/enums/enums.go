package enum

import (
	"html/template"

	"github.com/clickyab/services/codegen/annotate"

	"github.com/clickyab/services/codegen/plugins"

	"fmt"

	"bytes"
	"io/ioutil"

	"path/filepath"

	"github.com/goraz/humanize"
	"golang.org/x/tools/imports"
)

type enumPlugin struct {
}

type enumSingle struct {
	Type humanize.TypeName
	Vars []string
}

type enumCtx []enumSingle

const tmpl = `
// Code generated build with enum DO NOT EDIT.

package {{ .PackageName }}

{{ range $m := .Data }}
// IsValid try to validate enum value on ths type
func (e {{ $m.Type.Name }})IsValid() bool {
	return array.StringInArray(
		string(e),
		{{ range $v := $m.Vars }} string({{$v}}),
		{{ end }})
}

// Scan convert the json array ino string slice
func (e *{{ $m.Type.Name }}) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return trans.E("unsupported type")
	}
	if !{{ $m.Type.Name }}(b).IsValid() {
		return trans.E("invaid value")
	}
	*e = {{ $m.Type.Name }}(b)
	return nil
}

// Value try to get the string slice representation in database
func (e {{ $m.Type.Name }}) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, trans.E("invalid status")
	}
	return string(e), nil
}
{{ end }}
`

var tpl = template.Must(template.New("model").Parse(tmpl))

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (e enumPlugin) GetType() []string {
	return []string{"Enum"}
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (e enumPlugin) Finalize(c interface{}, p humanize.Package) error {
	var ctx enumCtx
	if c != nil {
		var ok bool
		ctx, ok = c.(enumCtx)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}
	out := &bytes.Buffer{}

	err := tpl.Execute(
		out,
		struct {
			PackageName string
			Data        enumCtx
		}{
			p.Name,
			ctx,
		})
	if err != nil {
		return err
	}
	path, _ := filepath.Split(p.Files[0].FileName)
	f := filepath.Join(path, "pkg_enums.gen.go")
	res, err := imports.Process(f, out.Bytes(), nil)
	if err != nil {
		fmt.Println(out.String())
		return err
	}
	if err := ioutil.WriteFile(f, res, 0644); err != nil {
		return err
	}

	return nil
}

// StructureIsSupported check for a function signature and if the function is supported in this
// interface
func (e enumPlugin) TypeIsSupported(f humanize.File, t humanize.TypeName) bool {
	return t.Type.GetDefinition() == "string"
}

// ProcessStructure the structure with its annotation. any error here means to stop the
// all process
// the first argument is the context. if its nil, means its the first run for this package.
// the result of this function is passed to the plugin next time for the next function
func (e enumPlugin) ProcessType(ctx interface{}, p humanize.Package, f humanize.File, t humanize.TypeName, a annotate.Annotate) (interface{}, error) {
	if ctx == nil {
		ctx = make(enumCtx, 0)
	}
	typ := enumSingle{
		Type: t,
		Vars: make([]string, 0),
	}
	for i := range p.Files {
		for _, c := range p.Files[i].Constants {
			if c.Type.GetDefinition() == t.Name {
				typ.Vars = append(typ.Vars, c.Name)
			}
		}
	}
	if len(typ.Vars) == 0 {
		return nil, fmt.Errorf("type %s has no variable", t.Name)
	}
	ctx = append(ctx.(enumCtx), typ)
	return ctx, nil
}

func (e enumPlugin) GetOrder() int {
	return 1000
}

func init() {
	plugins.Register(enumPlugin{})
}
