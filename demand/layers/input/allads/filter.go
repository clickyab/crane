package allads

import (
	"errors"

	"clickyab.com/crane/demand/filter/campaign"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/reducer"
	openrtb "clickyab.com/crane/openrtb/v2.5"
)

type mixer struct {
	f  []reducer.Filter
	fn func(int32, []string)
}

func (m mixer) Check(c entity.Context, a entity.Campaign) error {
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

	f = append(f, &campaign.Strategy{})
	if desktop {
		f = append(f, &campaign.Desktop{})
	}
	if os != "" {
		f = append(f, &campaign.OS{})
	}
	if whitelist {
		f = append(f, &campaign.WhiteList{})
	}
	if blacklist {
		f = append(f, &campaign.BlackList{})
	}
	if len(cat) > 0 {
		f = append(f, &campaign.Category{})
	}
	if province {
		f = append(f, &campaign.Province{})
	}
	if isp {
		f = append(f, &campaign.ISP{})
	}
	return f
}

func filterAppBuilder(province bool, latlon, carrier, appBrand string, isp, whitelist, blacklist bool,
	cat []openrtb.ContentCategory,
) []reducer.Filter {

	f := make([]reducer.Filter, 0)

	f = append(f, &campaign.Strategy{})
	if len(appBrand) > 0 {
		f = append(f, &campaign.AppBrand{})
	}
	if carrier != "" {
		f = append(f, &campaign.AppCarrier{})
	}
	if whitelist {
		f = append(f, &campaign.WhiteList{})
	}
	if blacklist {
		f = append(f, &campaign.BlackList{})
	}
	if len(cat) > 0 {
		f = append(f, &campaign.Category{})
	}
	if province {
		f = append(f, &campaign.Province{})
	}
	if isp {
		f = append(f, &campaign.ISP{})
	}
	if latlon != "" {
		f = append(f, &campaign.AreaInGlob{})
	}
	return f

}
