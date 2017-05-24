package gin

import (
	"bytes"

	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	"clickyab.com/exchange/services/codegen/annotate"
	"clickyab.com/exchange/services/codegen/plugins"

	"github.com/Sirupsen/logrus"
	"github.com/goraz/humanize"
	"golang.org/x/tools/imports"
)

type route struct {
	Route               string
	Method              string
	Function            string
	RoutePkg            string
	RouteMiddleware     []string
	RouteFuncMiddleware string
	RecType             string
	RecName             string
	Payload             string
	Resource            string
	Scope               string
}

type group struct {
	FileName            string
	GroupPkg            string
	StructName          string
	Group               string
	GroupMiddleware     []string
	GroupFuncMiddleware string
}

type context struct {
	G *group
	R *route
}

type finalRoute struct {
	group
	GroupRec string
	Full     bool
	Routes   []route
}

type routerPlugin struct {
}

const tmpl = `
// Code generated build with router DO NOT EDIT.

package {{ .GroupPkg }}

import (
	"gopkg.in/labstack/echo.v3"
)


// Routes return the route registered with this
func ({{ .GroupRec }} *{{ .StructName }}) Routes(r *echo.Echo, mountPoint string) {
	{{ if .Full }}
	groupMiddleware :=  []echo.MiddlewareFunc{
		{{ range $gm := .GroupMiddleware }}{{ $gm }},
		{{end}}
	}
	{{ if .GroupFuncMiddleware }}
	groupMiddleware = append(groupMiddleware, {{ $.GroupRec }}.{{ $.GroupFuncMiddleware|strip_type }}()...)
	{{ end }}
	group := r.Group(mountPoint + "{{ .Group }}", groupMiddleware...)

	{{ range $key ,$route := .Routes }}
	/* Route {{ $route | jsonize }} with key {{ $key }} */
	m{{ $key }} :=  []echo.MiddlewareFunc{
	{{ range $rm := .RouteMiddleware }}{{ $rm }},
		{{end}}
	}
	{{ if $route.Resource }}
	base.RegisterPermission("{{$route.Resource}}", "{{$route.Resource}}")
	m{{ $key }} = append(m{{ $key }}, authz.AuthorizeGenerator("{{$route.Resource}}",base.UserScope("{{$route.Scope}}"))){{ end }}
	{{ if $route.RouteFuncMiddleware }}
	m{{ $key }} = append(m{{ $key }}, {{ $.GroupRec }}.{{ $route.RouteFuncMiddleware|strip_type }}()...){{ end }}
	{{ if $route.Payload }} // Make sure payload is the last middleware
	m{{ $key }} = append(m{{ $key }}, middlewares.PayloadUnMarshallerGenerator({{$route.Payload}}{})){{ end }}
	group.{{ $route.Method }}("{{ $route.Route }}",{{ $.GroupRec }}.{{ $route.Function|strip_type }}, m{{ $key }}...)
	// End route with key {{ $key }}
	{{ end }}
	{{ end }}
	initializer.DoInitialize({{ .GroupRec }})
}
`

var (
	echoImportPath = "gopkg.in/labstack/echo.v3"
	validMethod    = []string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE", "CONNECT", "TRACE"}
	fMap           = template.FuncMap{
		"ucfirst":    ucFirst,
		"md5":        md5Sum,
		"strip_type": stripType,
		"jsonize":    jsonize,
	}
	tpl        = template.Must(template.New("gin").Funcs(fMap).Parse(tmpl))
	controller = regexp.MustCompile("Controller$")
)

func jsonize(in interface{}) string {
	res, _ := json.MarshalIndent(in, "\t", "\t")

	return string(res)
}

func ucFirst(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func md5Sum(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
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
func inArray(n string, h ...string) bool {
	for i := range h {
		if n == h[i] {
			return true
		}
	}

	return false
}

func removeStar(t humanize.Type) humanize.Type {
	if s, ok := t.(*humanize.StarType); ok {
		return s.Target
	}

	return t
}

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (r *routerPlugin) GetType() []string {
	return []string{"Route"}
}

// ProcessFunction the function with its annotation. any error here means to stop the
// all process
// the first argument is the context. if its nil, means its the first run for this package.
// the result of this function is passed to the plugin next time for the next function
func (r *routerPlugin) ProcessFunction(
	c interface{},
	pkg humanize.Package,
	p humanize.File,
	f humanize.Function,
	a annotate.Annotate,
) (interface{}, error) {
	var ctx []context
	if c != nil {
		var ok bool
		ctx, ok = c.([]context)
		if !ok {
			return nil, fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	if _, ok := a.Items["ignore"]; ok {
		// there is an ignore key
		return ctx, nil
	}

	data := route{}
	var ok bool
	data.Route, ok = a.Items["url"]
	if !ok {
		return returnErr("@Route.url")
	}

	data.Method, ok = a.Items["method"]
	if !ok {
		return returnErr("@Route.method")
	}

	data.Method = strings.ToUpper(data.Method)
	if !inArray(data.Method, validMethod...) {
		return nil, fmt.Errorf("invalid method %s", data.Method)
	}
	data.Function = f.Name
	data.RoutePkg = p.PackageName
	tmp := a.Items["middleware"]
	for _, m := range strings.Split(tmp, ",") {
		if mt := strings.Trim(m, "\t "); mt != "" {
			data.RouteMiddleware = append(data.RouteMiddleware, mt)
		}
	}

	payload, ok := a.Items["payload"]
	if ok {
		payload = strings.Trim(payload, " \n\t")
		t, e := pkg.FindType(payload)
		if e != nil {
			return nil, fmt.Errorf("can not find type %s as payload", payload)
		}
		if _, ok = t.Type.(*humanize.StructType); !ok {
			return nil, fmt.Errorf("type %s is not a structure, must be a structure for payload", payload)
		}
		data.Payload = payload
	}

	resource, ok := a.Items["resource"]
	if ok {
		resource = strings.Trim(resource, " \n\t")
		scopes := strings.Split(resource, ":")
		scope := "global"
		if len(scopes) == 2 {
			scope = scopes[1]
			resource = scopes[0]
		}
		// For routes with the resource, must add the authenticate middleware
		if !stringInArray("authz.Authenticate", data.RouteMiddleware...) {
			data.RouteMiddleware = append(data.RouteMiddleware, "authz.Authenticate")
		}
		data.Resource = resource
		data.Scope = scope
	}

	data.RecType = removeStar(f.Receiver.Type).GetDefinition()
	data.RecName = f.Receiver.Name
	// Maybe it have a middleware function?
	// find the function with the exact name
	fn := f.Name + "Middleware"
	if file, fnDef, ok := findFunction(pkg, fn); ok && isValidMiddleware(*file, *fnDef) {
		data.RouteFuncMiddleware = fn
	}

	n := context{
		G: nil,
		R: &data,
	}

	ctx = append(ctx, n)

	return ctx, nil
}

func (r *routerPlugin) ProcessStructure(
	c interface{},
	pkg humanize.Package,
	p humanize.File,
	f humanize.TypeName,
	a annotate.Annotate,
) (interface{}, error) {
	var ctx []context
	if c != nil {
		var ok bool
		ctx, ok = c.([]context)
		if !ok {
			return nil, fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}
	if _, ok := a.Items["ignore"]; ok {
		// there is an ignore key
		return ctx, nil
	}
	data := group{}
	var ok bool
	data.Group, ok = a.Items["group"]
	if !ok {
		return returnErr("@Route.group")
	}

	data.StructName = f.Name
	data.FileName = p.FileName

	data.GroupPkg = p.PackageName
	tmp := a.Items["middleware"]
	for _, m := range strings.Split(tmp, ",") {
		if mt := strings.Trim(m, "\t "); mt != "" {
			data.GroupMiddleware = append(data.GroupMiddleware, mt)
		}
	}

	// Maybe it have a middleware function?
	// find the function with the exact name
	fn := f.Name + ".middleware"

	if file, fnDef, ok := findFunction(pkg, fn); ok && isValidMiddleware(*file, *fnDef) {
		data.GroupFuncMiddleware = fn
	}

	n := context{
		G: &data,
		R: nil,
	}

	ctx = append(ctx, n)

	return ctx, nil
}

func (r *routerPlugin) mix(ctx []context) map[string]finalRoute {
	var result = make(map[string]finalRoute)
	for i := range ctx {
		if ctx[i].G != nil {
			res := finalRoute{
				group:    *ctx[i].G,
				GroupRec: "ctrl",
			}
			for j := range ctx {
				// Check against all groups
				if ctx[j].R != nil && ctx[i].G.StructName == ctx[j].R.RecType {
					// Match!
					res.Routes = append(res.Routes, *ctx[j].R)
					res.GroupRec = ctx[j].R.RecName
				}
			}
			res.Full = len(res.Routes) > 0
			result[res.group.Group] = res
		}
	}

	return result
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (r *routerPlugin) Finalize(c interface{}, _ humanize.Package) error {
	var ctx []context
	if c != nil {
		var ok bool
		ctx, ok = c.([]context)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	nctx := r.mix(ctx)
	//x, _ := json.MarshalIndent(nctx, "", "\t")
	//fmt.Println(string(x))

	for typ := range nctx {
		out := &bytes.Buffer{}
		if err := tpl.Execute(out, nctx[typ]); err != nil {
			return err
		}

		f := nctx[typ].FileName
		f = f[:len(f)-3] + ".gen.go"
		res, err := imports.Process(f, out.Bytes(), nil)
		if err != nil {
			fmt.Println(out.String())
			return err
		}
		if err := ioutil.WriteFile(f, res, 0644); err != nil {
			return err
		}
	}

	return nil
}

// FindFunction try to find a package level variable
func findFunction(p humanize.Package, t string) (*humanize.File, *humanize.Function, bool) {
	for i := range p.Files {
		for j := range p.Files[i].Functions {
			if p.Files[i].Functions[j].Name == t {
				return p.Files[i], p.Files[i].Functions[j], true
			}
		}
	}
	return nil, nil, false
}

func isValidMiddleware(file humanize.File, fn humanize.Function) bool {
	gin, ok := findEchoImport(file)
	if !ok {
		return false
	}
	if len(fn.Type.Parameters) != 1 {
		return false
	}

	if len(fn.Type.Results) != 1 {
		return false
	}

	typ := fn.Type.Results[0].Type
	if typ.GetDefinition() != gin+".HandlerFunc" {
		return false
	}
	typ = fn.Type.Parameters[0].Type
	if typ.GetDefinition() != gin+".HandlerFunc" {
		return false
	}

	return true
}

func isMatched(gin string, f humanize.Function) bool {
	if len(f.Type.Results) != 1 {
		return false
	}

	if len(f.Type.Parameters) != 1 {
		return false
	}

	if f.Receiver == nil {
		return false
	}

	if f.Type.Parameters[0].Type.GetDefinition() != fmt.Sprintf("%s.Context", gin) {
		return false
	}

	if f.Type.Results[0].Type.GetDefinition() != "error" {
		return false
	}

	return true
}
func findEchoImport(f humanize.File) (string, bool) {
	for i := range f.Imports {
		if f.Imports[i].Path == echoImportPath {
			return "echo", true
		}
	}

	return "", false
}

func (r *routerPlugin) FunctionIsSupported(file humanize.File, fn humanize.Function) bool {
	str, b := findEchoImport(file)
	if !b {
		return false
	}

	return isMatched(str, fn)
}

func (r *routerPlugin) StructureIsSupported(file humanize.File, fn humanize.TypeName) bool {
	return controller.MatchString(fn.Name)
}

func (r *routerPlugin) GetOrder() int {
	return 999
}

func stringInArray(in string, arr ...string) bool {
	for i := range arr {
		if arr[i] == in {
			return true
		}
	}

	return false
}

func init() {
	plugins.Register(&routerPlugin{})
}
