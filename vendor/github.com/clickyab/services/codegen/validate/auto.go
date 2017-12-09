package validate

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/clickyab/services/codegen/annotate"
	"github.com/clickyab/services/codegen/plugins"

	"bytes"
	"io/ioutil"

	"path/filepath"

	"sort"

	"github.com/clickyab/services/assert"
	"github.com/goraz/humanize"
	"golang.org/x/tools/imports"
)

type validatePlugin struct {
}

type validate struct {
	pkg  humanize.Package
	file humanize.File
	ann  annotate.Annotate
	typ  humanize.TypeName

	Map  []fieldMap
	Rec  string
	Type string
}

type fieldMap struct {
	Name string
	Json string
	Err  string
}

type context []validate

func (c context) Len() int {
	return len(c)
}

func (c context) Less(i, j int) bool {
	return strings.Compare(c[i].Type, c[j].Type) < 0
}

func (c context) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

var (
	validateFunc = `
// Code generated build with variable DO NOT EDIT.

package {{ .PackageName }}
// AUTO GENERATED CODE. DO NOT EDIT!
import (
	"gopkg.in/go-playground/validator.v9"
	"github.com/clickyab/services/framework/middleware"
	"context"
	"net/http"
)
	{{ range $m := .Data }}
	func ({{ $m.Rec }} *{{ $m.Type }}) Validate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := func(in interface{}) error {
			if v, ok := in.(interface {
				ValidateExtra(ctx context.Context, w http.ResponseWriter, r *http.Request) error
			}); ok {
				return v.ValidateExtra(ctx, w, r)
			}
			return nil
		}({{ $m.Rec }})
		if err != nil {
			return err
		}
		errs :=  validator.New().Struct({{ $m.Rec }})
		if errs == nil {
			return nil
		}
		res := middleware.GroupError{}
		for _, i := range errs.(validator.ValidationErrors) {
			switch i.Field() { {{ range $f := $m.Map }}
				case "{{ $f.Name }}":
					res["{{ $f.Json }}"] = trans.E("{{ $f.Err }}")
			{{ end }}
				default :
					logrus.Panicf("the field %s is not translated", i)
			}
		}
		if len(res) >0 {
			return res
		}
		return nil
	}
	{{ end }}
	`

	tpl = template.Must(template.New("validate").Parse(validateFunc))
)

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (e validatePlugin) GetType() []string {
	return []string{"Validate"}
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (e validatePlugin) Finalize(c interface{}, p humanize.Package) error {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	buf := &bytes.Buffer{}
	sort.Sort(ctx)
	err := tpl.Execute(buf, struct {
		Data        context
		PackageName string
	}{
		Data:        ctx,
		PackageName: p.Name,
	})
	if err != nil {
		return err
	}
	f := filepath.Dir(p.Files[0].FileName)
	f = filepath.Join(f, "validators.gen.go")
	res, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		fmt.Println(buf.String())
		return err
	}

	err = ioutil.WriteFile(f, res, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *validatePlugin) ProcessStructure(
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

	dt := validate{
		pkg:  pkg,
		file: p,
		ann:  a,
		typ:  f,

		Type: f.Name,
		Rec:  "pl",
	}
	dt.Map = getAllValidateFields(f.Type)

bigLoop:
	for i := range pkg.Files {
		for _, fn := range pkg.Files[i].Functions {
			if fn.Receiver != nil {
				rec := fn.Receiver.Type
				if s, ok := rec.(*humanize.StarType); ok {
					rec = s.Target
				}
				if f.Name == rec.GetDefinition() {
					dt.Rec = fn.Receiver.Name
					break bigLoop
				}
			}
		}
	}

	ctx = append(ctx, dt)
	return ctx, nil
}

// getAllValidateFields get all validation fields from struct
func getAllValidateFields(f humanize.Type) []fieldMap {
	var res = []fieldMap{}
	if len(f.(*humanize.StructType).Embeds) > 0 {
		res = append(res, getValidateFieldsFromEmbeds(f.(*humanize.StructType).Embeds)...)
	}
	if len(f.(*humanize.StructType).Fields) > 0 {
		res = append(res, getValidateFieldFromStruct(f.(*humanize.StructType).Fields)...)
	}
	return res
}

// getValidateFieldsFromEmbeds get validation fields from embed struct
func getValidateFieldsFromEmbeds(e []*humanize.Embed) (res []fieldMap) {
	for _, val := range e {
		//struct from another package
		pack := val.Package()
		if v, ok := val.Type.(*humanize.SelectorType); ok {
			if vv, ok := v.Type.(*humanize.IdentType); ok {
				foundType, err := pack.FindType(vv.Ident)
				assert.Nil(err)
				res = append(res, getValidateFieldFromStruct(foundType.Type.(*humanize.StructType).Fields)...)
				if len(foundType.Type.(*humanize.StructType).Embeds) != 0 {
					res = append(res, getValidateFieldsFromEmbeds(foundType.Type.(*humanize.StructType).Embeds)...)
				}
			}
		} else {
			if i, ok := val.Type.(*humanize.IdentType); ok {
				foundType, err := pack.FindType(i.Ident)
				assert.Nil(err)
				res = append(res, getValidateFieldFromStruct(foundType.Type.(*humanize.StructType).Fields)...)
				if len(foundType.Type.(*humanize.StructType).Embeds) != 0 {
					res = append(res, getValidateFieldsFromEmbeds(foundType.Type.(*humanize.StructType).Embeds)...)
				}
			}
		}
	}
	return
}

// getValidateFieldFromStruct get validation fields from normal struct or array struct
//TODO : not support pointers yet (we should decide if we want to support them)
func getValidateFieldFromStruct(f []*humanize.Field) []fieldMap {
	var res = []fieldMap{}
	for _, l := range f {
		if t, ok := l.Type.(*humanize.ArrayType); ok {
			if g, ok := t.Type.(*humanize.StructType); ok {
				res = append(res, getValidateFieldFromStruct(g.Fields)...)
			}
		}
		if l.Tags.Get("validate") != "" {
			t := fieldMap{
				Name: l.Name,
				Json: l.Tags.Get("json"),
				Err:  l.Tags.Get("error"),
			}

			if t.Json == "" {
				t.Json = t.Name
			}

			if t.Err == "" {
				t.Err = "invalid value"
			}

			res = append(res, t)
		}
	}
	return res
}

func (r *validatePlugin) StructureIsSupported(file humanize.File, fn humanize.TypeName) bool {
	return true
}

func (r *validatePlugin) GetOrder() int {
	return 5999
}

func init() {
	plugins.Register(&validatePlugin{})
}
