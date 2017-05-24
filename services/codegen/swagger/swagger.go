package swagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/codegen/annotate"
	"clickyab.com/exchange/services/codegen/plugins"

	"github.com/Sirupsen/logrus"
	"github.com/goraz/humanize"
	"gopkg.in/yaml.v2"
)

const (
	ginImportPath = "gopkg.in/labstack/echo.v3"
)

type swaggerType map[string]interface{}

type uriParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type,omitempty" yaml:",omitempty"`
	Schema      swaggerType `json:"schema,omitempty" yaml:",omitempty"`
	Description string      `json:"description"`
	In          string      `json:"in"`
	Required    bool        `json:"required"`
}

var (
	transferList = map[string]swaggerType{
		"time.Time":    {"type": "string", "format": "dateTime"},
		"sql.NullTime": {"type": "string", "format": "dateTime"},
		"clickyab.com/exchange/services/postgres/models/common.NullTime":   {"type": "string", "format": "dateTime"},
		"clickyab.com/exchange/services/postgres/models/common.NullString": {"type": "string"},
		"modules/balance/acc.Money":                                        {"type": "integer"},
		"modules/user/aaa.UserStatus":                                      {"type": "string"},
	}
)

type apiBodyInner struct {
	Ref string `yaml:"$ref" json:"$ref"`
}

type apiBody struct {
	Description string       `json:"description"`
	Produce     string       `json:"produces,omitempty" yaml:"produces,omitempty"`
	Schema      apiBodyInner `json:"schema"`
}

type apiHeader struct {
	Description string `json:"description"`
}

type apiMethod struct {
	Description string                 `json:"description"`
	Responses   map[string]interface{} `json:"responses"`
	Parameters  []uriParameter         `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Tags        []string               `json:"tags"`
}

type apiGroup struct {
	Description string               `yaml:"-"`
	Route       string               `yaml:"-"`
	Parameters  []uriParameter       `yaml:"-"`
	Paths       map[string]apiMethod `yaml:",inline"`
	Tags        []string             `yaml:"-"`
	protected   bool
}

type context map[string]apiGroup

type swaggerGenerator struct {
	workDir string
	domain  string
}

var (
	typeCache  = make(map[string]swaggerType)
	parameters = regexp.MustCompile("/:([^/]+)")
	query      = regexp.MustCompile("^_([a-zA-Z0-9_]+)_$")
)

func ucFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (rg *swaggerGenerator) GetType() []string {
	return []string{"Route"}
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (rg *swaggerGenerator) Finalize(c interface{}, p humanize.Package) error {
	res, err := json.Marshal(typeCache)
	assert.Nil(err)
	pkg := filepath.Join(rg.workDir, strings.Replace(p.Path, "/", "_", -1)+"_definitions.json")
	if err = ioutil.WriteFile(pkg, res, 0644); err != nil {
		return err
	}

	res, err = json.Marshal(c)
	assert.Nil(err)
	pkg = filepath.Join(rg.workDir, strings.Replace(p.Path, "/", "_", -1)+"_entries.json")
	err = ioutil.WriteFile(pkg, res, 0644)
	if err != nil {
		return err
	}

	return rg.mix()
}

func loadPattern(pattern string) ([][]byte, error) {
	m, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	res := make([][]byte, len(m))
	for i := range m {
		res[i], err = ioutil.ReadFile(m[i])
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (rg *swaggerGenerator) convertAI(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i := range in {
		var tmp interface{}
		if dd, ok := in[i].(map[interface{}]interface{}); ok {
			tmp = rg.convertIS(dd)
		} else if dd, ok := in[i].(map[string]interface{}); ok {
			tmp = rg.convertSS(dd)
		} else if dd, ok := in[i].([]interface{}); ok {
			tmp = rg.convertAI(dd)
		} else {
			tmp = in[i]
		}

		res[i] = tmp
	}

	return res
}

func (rg *swaggerGenerator) convertSS(in map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for i := range in {
		var tmp interface{}
		if dd, ok := in[i].(map[interface{}]interface{}); ok {
			tmp = rg.convertIS(dd)
		} else if dd, ok := in[i].(map[string]interface{}); ok {
			tmp = rg.convertSS(dd)
		} else if dd, ok := in[i].([]interface{}); ok {
			tmp = rg.convertAI(dd)
		} else {
			tmp = in[i]
		}

		res[i] = tmp
	}

	return res
}

func (rg *swaggerGenerator) convertIS(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for i := range in {
		var tmp interface{}
		if dd, ok := in[i].(map[interface{}]interface{}); ok {
			tmp = rg.convertIS(dd)
		} else if dd, ok := in[i].(map[string]interface{}); ok {
			tmp = rg.convertSS(dd)
		} else if dd, ok := in[i].([]interface{}); ok {
			tmp = rg.convertAI(dd)
		} else {
			tmp = in[i]
		}

		switch ii := i.(type) {
		case string:
			res[ii] = tmp
		default:
			panic(fmt.Sprintf("%[1]T => %[1]v", i))
		}
	}

	return res
}

func (rg *swaggerGenerator) mix() error {
	defs, err := loadPattern(filepath.Join(rg.workDir, "*_definitions.json"))
	if err != nil {
		return err
	}
	nttFinal := make(map[string]swaggerType)
	for i := range defs {
		tmp := make(map[string]swaggerType)
		err = json.Unmarshal(defs[i], &tmp)
		if err != nil {
			return err
		}

		for i := range tmp {
			nttFinal[i] = tmp[i]
		}
	}

	ntt, err := loadPattern(filepath.Join(rg.workDir, "*_entries.json"))
	if err != nil {
		return err
	}

	defsFinal := make(map[string]apiGroup)
	for i := range ntt {
		tmp := make(map[string]apiGroup)
		err = json.Unmarshal(ntt[i], &tmp)
		if err != nil {
			return err
		}

		for i := range tmp {
			defsFinal[i] = tmp[i]
		}
	}

	swagger := struct {
		Swagger string
		Info    struct {
			Version     string
			Title       string
			Description string
		}
		Host        string
		BasePath    string `yaml:"basePath"`
		Schemes     []string
		Consumes    []string
		Produces    []string
		Paths       map[string]apiGroup
		Definitions map[string]swaggerType
	}{
		Swagger: "2.0",
		Info: struct {
			Version     string
			Title       string
			Description string
		}{
			Version:     "1.0.0",
			Title:       "The Malooch API",
			Description: "Auto genertaed Malooch API",
		},
		Host:        rg.domain,
		BasePath:    "/api",
		Schemes:     []string{"http"},
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Paths:       defsFinal,
		Definitions: nttFinal,
	}

	data, err := yaml.Marshal(swagger)
	if err != nil {
		return err
	}

	tmp := make(map[string]interface{})
	err = yaml.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	res := rg.convertSS(tmp)

	jsonData, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(rg.workDir, "out.json"), jsonData, 0644)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(rg.workDir, "out.yaml"), data, 0644)
}

// FunctionIsSupported check for a function signature and if the function is supported in this
// interface
func (rg *swaggerGenerator) FunctionIsSupported(f humanize.File, fn humanize.Function) bool {
	str, b := findGinImport(f)
	if !b {
		return false
	}

	return isMatched(str, fn)
}

// ProcessFunction the function with its annotation. any error here means to stop the
// all process
// the first argument is the context. if its nil, means its the first run for this package.
// the result of this function is passed to the plugin next time for the next function
func (rg *swaggerGenerator) ProcessFunction(ctx interface{}, pkg humanize.Package, f humanize.File, fn humanize.Function, ann annotate.Annotate) (interface{}, error) {
	var c context
	var ok bool
	if c, ok = ctx.(context); !ok {
		c = make(context)
	}

	if fn.Receiver == nil {
		return nil, fmt.Errorf("reciever is nil")
	}

	t := fn.Receiver.Type
	if _, ok := t.(*humanize.StarType); ok {
		t = t.(*humanize.StarType).Target
	}

	var route *annotate.Annotate
	var docs humanize.Docs
	switch tt := t.(type) {
	case *humanize.IdentType:
		st, err := pkg.FindType(tt.Ident)
		if err != nil {
			return nil, err
		}
		docs = st.Docs
		gann, err := annotate.LoadFromComment(strings.Join(docs, "\n"))
		if err != nil {
			return nil, err
		}

		for i := range gann {
			if gann[i].Name == "Route" {
				route = &gann[i]
				break
			}
		}

	default:
		logrus.Panicf("TODO : %T is not supported", t)
	}

	if route == nil {
		return nil, fmt.Errorf("the group has no annotate")
	}

	groupName, ok := route.Items["group"]
	if !ok {
		return nil, fmt.Errorf("no @Route.group")
	}
	var pp []string
	groupName, pp = changeUrl(groupName)
	urlName, pp2 := changeUrl(ann.Items["url"])
	final := groupName + urlName
	ag, ok := c[final]
	if !ok {
		ag = apiGroup{
			Description: getDoc(docs),
			Paths:       make(map[string]apiMethod),
		}

		mid, ok := route.Items["middleware"]
		if ok {
			midA := strings.Split(mid, ",")
			for _, mm := range midA {
				if strings.Trim(mm, " \t\n") == "authz.Authenticate" {
					ag.protected = true
				}
			}
		} else if _, ok := route.Items["resource"]; ok {
			ag.protected = true
		}
		ag.Route = groupName
		for _, p := range pp {
			if v, ok := route.Items["_"+p+"_"]; ok {
				delete(route.Items, "_"+p+"_")
				parts := strings.Split(v, ",")
				if len(parts) == 2 {
					ag.Parameters = append(ag.Parameters,
						uriParameter{
							Name:        p,
							Type:        parts[0],
							Description: parts[1],
							In:          "path",
							Required:    true,
						},
					)
				} else if len(parts) == 1 {
					ag.Parameters = append(ag.Parameters,
						uriParameter{
							Name:     p,
							Type:     parts[0],
							In:       "path",
							Required: true,
						},
					)
				}
			} else {
				ag.Parameters = append(ag.Parameters,
					uriParameter{
						Name:     p,
						Type:     "string",
						In:       "path",
						Required: true,
					},
				)
			}
		}

		for i, v := range route.Items {
			qParam := query.FindStringSubmatch(i)
			if len(qParam) > 0 {
				parts := strings.Split(v, ",")
				if len(parts) == 2 {
					ag.Parameters = append(ag.Parameters,
						uriParameter{
							Name:        qParam[1],
							Type:        parts[0],
							Description: parts[1],
							In:          "query",
							Required:    false,
						},
					)
				} else if len(parts) == 1 {
					ag.Parameters = append(ag.Parameters,
						uriParameter{
							Name:     qParam[1],
							Type:     parts[0],
							In:       "query",
							Required: false,
						},
					)
				} else {
					ag.Parameters = append(ag.Parameters,
						uriParameter{
							Name:     qParam[1],
							Type:     "string",
							In:       "query",
							Required: false,
						},
					)
				}
			}
		}

		ag.Tags = append(ag.Tags, pkg.Name)
		tags := strings.Trim(route.Items["tags"], " \n\t")
		if tags != "" {
			ag.Tags = append(ag.Tags, strings.Split(tags, ",")...)
		}
	}
	method := ann.Items["method"]
	_, ok = ag.Paths[method]
	if ok {
		return nil, fmt.Errorf("duplicate method")
	}
	am := apiMethod{
		Description: getDoc(fn.Docs),
		Responses:   make(map[string]interface{}),
		Parameters:  ag.Parameters,
	}
	for _, p := range pp2 {
		if v, ok := ann.Items["_"+p+"_"]; ok {
			// ok remove it
			delete(ann.Items, "_"+p+"_")
			parts := strings.Split(v, ",")
			if len(parts) == 2 {
				am.Parameters = append(
					am.Parameters,
					uriParameter{
						Name:        p,
						Type:        strings.Trim(parts[0], " \t"),
						Description: parts[1],
						In:          "path",
						Required:    true,
					},
				)

			} else if len(parts) == 1 {
				am.Parameters = append(
					am.Parameters,
					uriParameter{
						Name:     p,
						Type:     strings.Trim(parts[0], " \t"),
						In:       "path",
						Required: true,
					},
				)
			}
		} else {
			am.Parameters = append(
				am.Parameters,
				uriParameter{
					Name:     p,
					Type:     "string",
					In:       "path",
					Required: true,
				},
			)

		}
	}
	for i, v := range ann.Items {
		qParam := query.FindStringSubmatch(i)
		if len(qParam) > 0 {
			parts := strings.Split(v, ",")
			if len(parts) == 2 {
				am.Parameters = append(am.Parameters,
					uriParameter{
						Name:        qParam[1],
						Type:        strings.Trim(parts[0], " \t"),
						Description: parts[1],
						In:          "query",
						Required:    false,
					},
				)
			} else if len(parts) == 1 {
				am.Parameters = append(am.Parameters,
					uriParameter{
						Name:     qParam[1],
						Type:     strings.Trim(parts[0], " \t"),
						In:       "query",
						Required: false,
					},
				)
			} else {
				am.Parameters = append(am.Parameters,
					uriParameter{
						Name:     qParam[1],
						Type:     "string",
						In:       "query",
						Required: false,
					},
				)
			}
		}
	}
	pro := ag.protected
	if !pro {
		mid, ok := ann.Items["middleware"]
		if ok {
			midA := strings.Split(mid, ",")
			for _, mm := range midA {
				if strings.Trim(mm, " \t\n") == "authz.Authenticate" {
					pro = true
				}
			}
		} else if _, ok := ann.Items["resource"]; ok {
			pro = true
		}
	}

	if pro {
		// Currently add token parameter
		am.Parameters = append(
			am.Parameters,
			uriParameter{
				Name:        "token",
				Type:        "string",
				In:          "header",
				Description: "the security token, get it from login route",
				Required:    true,
			},
		)

		am.Responses[fmt.Sprintf("%d", 403)] = map[string]interface{}{
			"description": "forbidden, you have no access here",
			"schema": map[string]interface{}{
				"title": "forbidden",
				"type":  "object",
				"properties": map[string]interface{}{
					"error": map[string]interface{}{
						"properties": map[string]interface{}{
							"params": map[string]interface{}{
								"items": map[string]string{
									"type": "string",
								},
								"type": "array",
							},
							"text": map[string]string{
								"type": "string",
							},
						},
						"type": "object",
					},
				},
			},
		}

		am.Responses[fmt.Sprintf("%d", 401)] = map[string]interface{}{
			"description": "you are not authorized",
			"schema": map[string]interface{}{
				"title": "not_authorized",
				"type":  "object",
				"properties": map[string]interface{}{
					"error": map[string]interface{}{
						"properties": map[string]interface{}{
							"params": map[string]interface{}{
								"items": map[string]string{
									"type": "string",
								},
								"type": "array",
							},
							"text": map[string]string{
								"type": "string",
							},
						},
						"type": "object",
					},
				},
			},
		}
	}
	payload, ok := ann.Items["payload"]
	if !ok {
		payload, ok = ann.Items["meta_payload"]
	}
	if ok {
		// there is a payload, load type
		n, _, err := findType(pkg, payload)
		if err != nil {
			return nil, err
		}
		am.Parameters = append(
			am.Parameters,
			uriParameter{
				Name:   "payload_data",
				In:     "body",
				Schema: swaggerType{"$ref": "#/definitions/" + n},
			},
		)
	}

	for i, v := range ann.Items {
		code, err := strconv.ParseInt(i, 10, 0)
		if err != nil {
			continue
		}

		n, _, err := findType(pkg, v)
		if err != nil {
			return nil, err
		}
		var (
			produce string
		)
		if code == 200 {
			produce = ann.Items["produces"]
		}
		am.Responses[fmt.Sprintf("%d", code)] = apiBody{
			Schema: apiBodyInner{
				Ref: "#/definitions/" + n,
			},
			Produce: produce,
		}

	}
	am.Tags = append(am.Tags, ag.Tags...)
	tags := strings.Trim(ann.Items["tags"], " \n\t")
	if tags != "" {
		am.Tags = append(am.Tags, strings.Split(tags, ",")...)
	}

	ag.Paths[method] = am
	c[final] = ag
	return c, nil
}

func getDoc(d humanize.Docs) string {
	var res []string
	for _, s := range d {
		s := strings.Trim(s, " /")
		if s[0] == '@' {
			break
		}
		res = append(res, s)
	}

	return strings.Join(res, " ")
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

func findGinImport(f humanize.File) (string, bool) {
	for i := range f.Imports {
		if f.Imports[i].Path == ginImportPath {
			// TODO : make sure edit this if you change this
			return "echo", true
		}
	}

	return "", false
}

func changeUrl(s string) (string, []string) {
	var res []string
	match := parameters.FindAllStringSubmatch(s, -1)
	for i := range match {
		res = append(res, match[i][1])
	}
	s = parameters.ReplaceAllString(s, "/{${1}}")

	return s, res
}

func findType(pkg humanize.Package, t string) (string, swaggerType, error) {
	theName := pkg.Name + "_" + t

	if res, ok := typeCache[theName]; ok {
		return theName, res, nil
	}

	p, err := pkg.FindType(t)
	if err != nil {
		// there is special case with external package
		tmp := strings.Split(t, ".")
		if len(tmp) == 2 {
			im, er := pkg.FindImport(tmp[0])
			if er != nil {
				return "", nil, er
			}
			newPkg, er := humanize.ParsePackage(im.Path)
			if er != nil {
				return "", nil, er
			}
			return findType(*newPkg, tmp[1])
		}
		return "", nil, err
	}

	res, err := goToRaml(pkg, p.Type)
	if err != nil {
		return "", nil, err
	}
	typeCache[theName] = res
	return theName, res, nil
}

func goToRaml(pkg humanize.Package, tn humanize.Type) (swaggerType, error) {

	switch t := tn.(type) {
	case *humanize.StarType:
		return goToRaml(pkg, t.Target)
	case *humanize.StructType:
		return structToRaml(pkg, t)
	case *humanize.IdentType:
		return identToRaml(pkg, t)
	case *humanize.SelectorType:
		return selectorType(pkg, t)
	case *humanize.ArrayType:
		return sliceToRml(pkg, t)
	case *humanize.MapType:
		return mapToRaml(pkg, t)
	default:
		logrus.Panicf("TODO  %T => %+v , %s", tn, tn, tn.GetDefinition())
	}

	return nil, fmt.Errorf("format is not supported")
}

func selectorType(pkg humanize.Package, i *humanize.SelectorType) (swaggerType, error) {
	for param, res := range transferList {
		ident := strings.Split(param, ".")
		assert.True(len(ident) == 2, "[BUG] len is not 2")
		if i.Type.(*humanize.IdentType).Ident == ident[1] && i.Package().Path == ident[0] {
			return res, nil
		}
	}

	newPkg := i.Package()
	typ, err := newPkg.FindType(i.Type.(*humanize.IdentType).Ident)
	if err != nil {
		return nil, err
	}

	return goToRaml(*newPkg, typ.Type)
}

func identToRaml(pkg humanize.Package, i *humanize.IdentType) (swaggerType, error) {
	for param, res := range transferList {
		ident := strings.Split(param, ".")
		assert.True(len(ident) == 2, "[BUG] len is not 2")
		if i.Ident == ident[1] && pkg.Path == ident[0] {
			return res, nil
		}
	}
	if i.Ident == "string" || i.Ident == "interface" {
		return swaggerType{"type": "string"}, nil
	} else if i.Ident == "int64" || i.Ident == "int" {
		return swaggerType{"type": "integer"}, nil
	} else if i.Ident == "float64" || i.Ident == "float32" {
		return swaggerType{"type": "number"}, nil
	} else if i.Ident == "bool" {
		return swaggerType{"type": "boolean"}, nil
	}
	_, typ, err := findType(pkg, i.Ident)
	if err != nil {
		return nil, err
	}

	return typ, nil
}

func sliceToRml(pkg humanize.Package, arr *humanize.ArrayType) (swaggerType, error) {
	var (
		theName = swaggerType{"type": "string"}
		err     error
	)
	switch t := arr.Type.(type) {
	case *humanize.IdentType:
		theName, err = identToRaml(pkg, t)
		if err != nil {
			return nil, err
		}
	case *humanize.SelectorType:
		i, err := pkg.FindImport(t.Package().Path)
		if err != nil {
			return nil, err
		}

		newPkg, err := humanize.ParsePackage(i.Path)
		if err != nil {
			return nil, err
		}
		theName, err = identToRaml(*newPkg, t.Type.(*humanize.IdentType))
		if err != nil {
			return nil, err
		}
	case *humanize.MapType:
		theName, err = mapToRaml(pkg, t)
		if err != nil {
			return nil, err
		}
	default:
		//fmt.Printf("%T", t)
	}

	return swaggerType{"type": "array", "items": theName}, nil
}

func structToRaml(pkg humanize.Package, st *humanize.StructType) (swaggerType, error) {
	var err error
	res := make(swaggerType)
	res["type"] = "object"
	props := make(swaggerType)
	for _, f := range st.Fields {
		tags := strings.Split(f.Tags.Get("json"), ",")
		name := strings.Trim(tags[0], " ")
		if name == "-" {
			continue
		}

		if ucFirst(f.Name) != f.Name {
			continue
		}

		if name == "" {
			name = f.Name
		}

		props[name], err = goToRaml(pkg, f.Type)
		if err != nil {
			return nil, err
		}
	}

	for _, e := range st.Embeds {
		tags := strings.Split(e.Tags.Get("json"), ",")
		name := strings.Trim(tags[0], " ")
		if name == "-" {
			continue
		}
		var (
			nameSt string
			raml   swaggerType
		)
		switch st := e.Type.(type) {
		case *humanize.IdentType:
			nameSt, raml, err = findType(pkg, st.Ident)
		case *humanize.SelectorType:
			nameSt = strings.Replace(st.GetDefinition(), ".", "_", -1)
			raml, err = selectorType(pkg, st)
		default:
			logrus.Panicf("TODO %T => %+v", e.Type, e.Type)
		}

		if err != nil {
			return nil, err
		}
		if name == "" {
			// Embeded with no name
			for k, v := range raml["properties"].(swaggerType) {
				props[k] = v
			}
		} else {
			props["name"] = nameSt
		}

	}

	if len(props) > 0 {
		res["properties"] = props
	} else {
		res["properties"] = make(map[string]interface{})
	}
	return res, nil
}

func mapToRaml(pkg humanize.Package, st *humanize.MapType) (swaggerType, error) {
	var err error
	res := make(swaggerType)
	res["type"] = "object"
	props := make(swaggerType)
	if _, ok := st.Value.(*humanize.InterfaceType); ok {
		props["type"] = "string"
	} else if arr, ok := st.Value.(*humanize.ArrayType); ok {
		props, _ = sliceToRml(pkg, arr)
	}

	if len(props) == 0 {
		props["type"] = "array"
		props["items"], err = goToRaml(pkg, st.Value)
		if err != nil {
			return nil, err
		}

	}
	res["additionalProperties"] = props

	return res, nil
}

func (rg *swaggerGenerator) GetOrder() int {
	return 5000
}

func init() {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "127.0.0.1"
	}
	wd := os.Getenv("WORK_DIR")
	if wd == "" {
		wd = os.TempDir()
	}
	workDir := filepath.Join(wd, "swagger")
	assert.Nil(os.MkdirAll(workDir, 0744))
	plugins.Register(&swaggerGenerator{workDir: workDir, domain: domain})
}
