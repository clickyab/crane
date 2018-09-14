package allads

import (
	"errors"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/filter"
	"clickyab.com/crane/demand/reducer"
	"clickyab.com/crane/openrtb"
)

type mixer struct {
	f  []reducer.Filter
	fn func(int32, []string)
}

func (m mixer) Check(c entity.Context, a entity.Creative) error {
	es := make([]string, 0)
	for i := range m.f {
		if e := m.f[i].Check(c, a); e != nil {
			es = append(es, e.Error())
		}
	}
	if len(es) > 0 {
		m.fn(a.ID(), es)
		return errors.New("some filter didn't pass")
	}
	return nil
}

// Mix try to mix multiple filter to single function so there is no need to
// call Apply more than once
func Mix(fn func(adid int32, err []string), f ...reducer.Filter) reducer.Filter {
	return &mixer{f: f, fn: fn}
}

func filterWebBuilder(desktop, province bool, os string, isp,
	whitelist, blacklist bool, cat []openrtb.ContentCategory) []reducer.Filter {
	f := make([]reducer.Filter, 0)

	f = append(f, &filter.Strategy{})
	if desktop {
		f = append(f, &filter.Desktop{})
	}
	if os != "" {
		f = append(f, &filter.OS{})
	}
	if whitelist {
		f = append(f, &filter.WhiteList{})
	}
	if blacklist {
		f = append(f, &filter.BlackList{})
	}
	if len(cat) > 0 {
		f = append(f, &filter.Category{})
	}
	if province {
		f = append(f, &filter.Province{})
	}
	if isp {
		f = append(f, &filter.ISP{})
	}
	return f
}

func filterAppBuilder(province bool, latlon, carrier, appBrand string, isp, whitelist, blacklist bool,
	cat []openrtb.ContentCategory,
) []reducer.Filter {

	f := make([]reducer.Filter, 0)

	f = append(f, &filter.Strategy{})
	if len(appBrand) > 0 {
		f = append(f, &filter.AppBrand{})
	}
	if carrier != "" {
		f = append(f, &filter.AppCarrier{})
	}
	if whitelist {
		f = append(f, &filter.WhiteList{})
	}
	if blacklist {
		f = append(f, &filter.BlackList{})
	}
	if len(cat) > 0 {
		f = append(f, &filter.Category{})
	}
	if province {
		f = append(f, &filter.Province{})
	}
	if isp {
		f = append(f, &filter.ISP{})
	}
	if latlon != "" {
		f = append(f, &filter.AreaInGlob{})
	}
	return f

}
