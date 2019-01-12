package native

import (
	"fmt"
	"html/template"
	"sync"

	"github.com/clickyab/services/assert"
)

const (
	normalTemplate = `
{{range $index, $results := .}}
	<div>
        <a href='{{.Click}}' target='_blank'>
            <img src='{{.Image}}'/‎>
            <img style="display: none;" src='{{.Impression}}'/‎>
            <div>{{.Title}}</div>
        </a>
    </div>
{{end}}
`

	textTemplate = `
{{range $index, $results := .}}
	<div>
        <a href='{{.Click}}' target='_blank'>
            <img style="display: none;" src='{{.Impression}}'/‎>
            <div>{{.Title}}</div>
        </a>
    </div>
{{end}}
`
)

var (
	all  = make(map[string]*nativeTemplate)
	lock = sync.RWMutex{}

	templateFuncs = template.FuncMap{
		"isEven": isEven,
		"isOdd":  isOdd,
	}
	normalTpl = template.Must(template.New("native-normal").Funcs(templateFuncs).Parse(normalTemplate))
	textTpl   = template.Must(template.New("native-text").Funcs(templateFuncs).Parse(textTemplate))
)

type nativeTemplate struct {
	Counts   []int
	Image    bool
	Template *template.Template
}

func isEven(x int) bool {
	return (x+1)%2 == 0
}

func isOdd(x int) bool {
	return (x+1)%2 != 0
}

func registerNativeTemplate(name string, template *template.Template, image bool, counts ...int) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := all[name]
	assert.False(ok, "template name is already here")

	all[name] = &nativeTemplate{
		Counts:   counts,
		Image:    image,
		Template: template,
	}
}

func getNativeTemplate(name string) (*nativeTemplate, error) {
	lock.RLock()
	defer lock.RUnlock()

	n, ok := all[name]
	if !ok {
		return nil, fmt.Errorf("template with name %s is invalid", name)
	}

	return n, nil
}

func init() {
	registerNativeTemplate("grid3x", normalTpl, true, 3, 6, 12)
	registerNativeTemplate("grid4x", normalTpl, true, 2, 4, 8, 12)
	registerNativeTemplate("single", normalTpl, true, 1)
	registerNativeTemplate("text-list", textTpl, false, 1, 3, 4, 6, 8, 12)
	registerNativeTemplate("vertical", normalTpl, true, 1, 3, 4, 6, 8, 12)
}
