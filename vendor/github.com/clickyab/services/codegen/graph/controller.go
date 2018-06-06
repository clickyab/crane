package graph

const ctrl = `// Code generated build with graph DO NOT EDIT.

package {{ .ControllerPackageName }}

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"clickyab.com/crab/modules/domain/middleware/domain"
	"clickyab.com/crab/modules/user/middleware/authz"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/permission"
	"github.com/rs/xmux"
	"{{ .Model }}"
)

var (
	maxRange{{ .Data.Entity|ucfirst }} = config.RegisterDuration("srv.graph.max_range", 24*90*time.Hour, "maximum possible for graph date range")
	epoch{{ .Data.Entity|ucfirst }}    time.Time
	layout{{ .Data.Entity|ucfirst }}   = "2006010215"
	factor{{ .Data.Entity|ucfirst }} time.Duration = {{ if eq .Data.Period "hourly" }}1{{ else }}24{{ end }}

)

type graph{{ .Data.Entity|ucfirst }}Response struct {
	Format string      ` + "`json:\"format\"`" + `
	From   time.Time   ` + "`json:\"from\"`" + `
	To     time.Time   ` + "`json:\"to\"`" + `
	Type   string      ` + "`json:\"type\"`" + `
	Data   []graph{{ .Data.Entity|ucfirst }}Data ` + "`json:\"data\"`" + `
}

type graph{{ .Data.Entity|ucfirst }}Data struct {
	Title  string  ` + "`json:\"title\"`" + `
	Name   string  ` + "`json:\"name\"`" + `
	Hidden bool    ` + "`json:\"hidden\"`" + `
	Type   string  ` + "`json:\"type\"`" + `
	Order   int64      ` + "`json:\"order\"`" + `
	Data   []float64 ` + "`json:\"data\"`" + `
	Sum float64 ` + "`json:\"sum\"`" + `
	Avg float64 ` + "`json:\"avg\"`" + `
	Min float64 ` + "`json:\"min\"`" + `
	Max float64 ` + "`json:\"max\"`" + `
	OmitEmpty bool ` + "`json:\"-\"`" + `
}

// @Route {
//		url = {{ .Data.URL }}
//		method = get
//		resource = {{ .Data.View.Total }}
//		_from_ = string , from date rfc3339 ex:2002-10-02T15:00:00.05Z
//		_to_ = string , to date rfc3339 ex:2002-10-02T15:00:00.05Z
//		200 = graph{{ .Data.Entity|ucfirst }}Response{{ range $f := .Data.Conditions }}{{ if $f.Filter }}
//		_{{ $f.Data }}_ = string , filter the {{ $f.Data }} field valid values are {{ $f.FilterValid }}{{ end }}{{ end }}{{ range $f := .Data.Conditions }}{{ if $f.Searchable }}
//		_{{ $f.Data }}_ = string , search the {{ $f.Data }} field {{ end }}{{ end }}
// }
func (ctrl *Controller) graph{{ .Data.Entity|ucfirst }}(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	usr := authz.MustGetUser(ctx)
	domain := domain.MustGetDomain(ctx)
	filter := make(map[string]string)
	{{ range $f := .Data.Conditions }}
	{{ if $f.Filter }}if e := r.URL.Query().Get("{{ $f.Data }}"); e != "" && {{ $.PackageName }}.{{ $f.FieldTypeString }}(e).IsValid() {
		filter["{{ if ne $f.Transform "" }}{{ $f.Transform }}{{else}}{{ $f.Data }}{{end}}"] = e
	}{{ end }}{{ end }}
	search := make(map[string]string)
	{{ range $f := .Data.Conditions }}
	{{ if $f.Searchable }}if e := r.URL.Query().Get("{{ $f.Data }}"); e != "" {
		search["{{ if ne $f.Transform "" }}{{ $f.Transform }}{{else}}{{ $f.Data }}{{end}}"] = e
	}{{ end }}	{{ end }}
	params := make(map[string]string)
	for _, i := range xmux.Params(ctx) {
		params[i.Name] = xmux.Param(ctx,i.Name)
	}
	from, to, err := dateParam{{ .Data.Entity|ucfirst }}(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	if err != nil {
		ctrl.BadResponse(w, err)
		return
	}
	l, fn := dateRange{{ .Data.Entity|ucfirst }}(from, to)
	tm := make(map[string]graph{{ .Data.Entity|ucfirst }}Data)
	pc := permission.NewInterfaceComplete(usr, usr.ID, "{{ .Data.View.Perm }}", "{{ .Data.View.Scope }}",domain.ID)
	{{ range $i, $f := .Data.Layouts }}

	tm["{{ $f.Name }}"] = graph{{ $.Data.Entity|ucfirst }}Data{
					Name:   "{{ $f.Name }}",
					Title:  "{{ $f.Title }}",
					Type:   "{{ $f.Type }}",
					Hidden: {{ $f.Hidden }},
					Order: {{ $f.Order }},
					OmitEmpty: {{ $f.OmitEmpty }},
					Data:   make([]float64, l),
				}{{ end }}
	for i, v := range {{ .PackageName }}.New{{ .PackageName|ucfirst }}Manager().{{ .Data.Fill }}(pc, filter, search, params, from, to) {
		m, err := fn(v.{{ .Data.Key }})
		assert.Nil(err)
		{{ range $i, $f := .Data.Layouts }}
		tx{{ $f.Name }}:= tm["{{ $f.Name }}"]
		cv{{ $f.Name }}:= {{ if ne $f.FieldType "float64" }}float64(v.{{ $f.Field }}){{ else }}v.{{ $f.Field }}{{ end }}
		tm["{{ $f.Name }}"].Data[m] += cv{{ $f.Name }}
		tx{{ $f.Name }}.Sum += cv{{ $f.Name }}
		if i == 0 {
			tx{{ $f.Name }}.Min = cv{{ $f.Name }}
			tx{{ $f.Name }}.Max = cv{{ $f.Name }}
		} else {
			if  cv{{ $f.Name }} > tx{{ $f.Name }}.Max {
				tx{{ $f.Name }}.Max = cv{{ $f.Name }}
			}
			if tx{{ $f.Name }}.Min >  cv{{ $f.Name }} {
				tx{{ $f.Name }}.Min = cv{{ $f.Name }}
			}
		}
		tm["{{ $f.Name }}"] = tx{{ $f.Name }}
		{{ end }}}
	res := graph{{ .Data.Entity|ucfirst }}Response{
		From:   from,
		To:     to,
		Format: "{{ .Data.Period }}", // hourly|daily|weekly|monthly|yearly
		Type: "{{ .Data.Scale }}", //  number|percent
		Data: make([]graph{{ .Data.Entity|ucfirst }}Data, 0),
	}
	for _, v := range tm {
		if v.Sum == 0 && v.OmitEmpty {
			continue
		}
		if l!=0{
			v.Avg = v.Sum / float64(l)
		}
		res.Data = append(res.Data, v)
	}
	ctrl.OKResponse(w, res)
}

func dateParam{{ .Data.Entity|ucfirst }}(f, t string) (time.Time, time.Time, error) {
	from, err := time.Parse(time.RFC3339Nano, f)
	from = from.Truncate(time.Hour * factor{{ .Data.Entity|ucfirst }})
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("wrong date format")
	}
	to, err := time.Parse(time.RFC3339Nano, t)
	to = to.Truncate(time.Hour * factor{{ .Data.Entity|ucfirst }})
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("wrong date format")
	}
	if to.Before(from) {
		from, to = to, from
	}
	if from.IsZero() && to.IsZero() {
		to = time.Now()
		from = to.AddDate(0, 0, -maxRange{{ .Data.Entity|ucfirst }}.Int())
	} else if from.IsZero() {
		from = to.AddDate(0, 0, -maxRange{{ .Data.Entity|ucfirst }}.Int())
	} else if to.IsZero() {
		to = from.AddDate(0, 0, maxRange{{ .Data.Entity|ucfirst }}.Int())
	}

	if to.After(time.Now()) {
		to = time.Now()
	}

	if from.Before(to.AddDate(0, 0, -maxRange{{ .Data.Entity|ucfirst }}.Int())) {
		from = to.AddDate(0, 0, -maxRange{{ .Data.Entity|ucfirst }}.Int())
	}

	if from.Before(epoch{{ .Data.Entity|ucfirst }}) {
		from = epoch{{ .Data.Entity|ucfirst }}
	}

	return from, to, nil
}

{{ if ne .Data.KeyType "time.Time" }}
func timeToID{{ .Data.Entity|ucfirst }}(d time.Time) int64 {
	h := int64(d.Truncate(time.Hour * factor{{ .Data.Entity|ucfirst }}).Sub(epoch{{ .Data.Entity|ucfirst }}).Hours())
	return {{ if eq .Data.Period "hourly" }}h + 1{{ else }}(h / 24) + 1{{ end }}
}
{{ end }}

func dateRange{{ .Data.Entity|ucfirst }}(f, t time.Time) (int, func({{ .Data.KeyType }}) (int, error)) {
	from := f.Truncate(time.Hour * factor{{ .Data.Entity|ucfirst }})
	to := t.Truncate(time.Hour * factor{{ .Data.Entity|ucfirst }})
	res := make(map[string]int)
	for i := 0; ; i++ {
		x := {{ if eq .Data.Period "hourly" }}from.Add(time.Hour * time.Duration(i)){{ else }}from.AddDate(0, 0, i){{ end }}
		if x.After(to) {
			break
		}
		res[{{ if eq .Data.KeyType "time.Time" }}x.Format(layout{{ .Data.Entity|ucfirst }}){{ else }}fmt.Sprint(timeToID{{ .Data.Entity|ucfirst }}(x)){{ end }}] = i
	}
	return len(res), func(m {{ .Data.KeyType }}) (int, error) {
		{{ if eq .Data.KeyType "time.Time" }}
		m = m.Truncate(time.Hour * factor{{ .Data.Entity|ucfirst }})
		if v, ok := res[m.Format(layout{{ .Data.Entity|ucfirst }})]; ok {
		{{ else }}
		if v, ok := res[fmt.Sprint(m)] ; ok {
		{{ end }}
			return v, nil
		}
		return 0, errors.New("out of range. probably mismatch key type. check {{ .Data.Fill }} annotation (e.g. daily to hourly or vice versa)")
	}
}

func init() {
	epoch{{ .Data.Entity|ucfirst }}, _ = time.Parse(layout{{ .Data.Entity|ucfirst }}, "{{ .Data.Epoch }}")
}

`
