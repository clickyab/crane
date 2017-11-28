package builder

import (
	"fmt"
	"net"
	"net/url"

	"net/http"

	"strings"

	"clickyab.com/gad/builder/cyos"
	"clickyab.com/gad/ip2location"
	"clickyab.com/gad/models"
	"github.com/clickyab/services/random"
	"github.com/mssola/user_agent"
)

// SetType is the type setter for context
func SetType(typ string) ShowOptionSetter {
	return func(options *Context) (*Context, error) {
		if typ != "vast" && typ != "App" && typ != "web" && typ != "native" {
			return nil, fmt.Errorf("type is not supported %s", typ)
		}
		options.common.Type = typ
		return options, nil
	}
}

// SetIP is the IP setter for context, also it extract the IP information
func SetIP(ip string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ipv4 := net.ParseIP(ip)
		if ipv4 == nil {
			return nil, fmt.Errorf("invalid IP %s", ip)
		}
		o.common.IP = ipv4

		// TODO : get extra data from this and add it to show context
		var l ip2location.LocationData
		o.common.ProvinceID, o.common.ISPID, l = ip2location.GetProvinceISPByIP(ipv4)
		o.common.Country, o.common.City, o.common.Province, o.common.Isp = l.Country, l.City, l.Province, l.ISP
		return o, nil
	}
}

// SetUserAgent try to set user agent and all related things
func SetUserAgent(ua string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		uaO := user_agent.New(ua)
		o.common.UserAgent = ua
		o.common.Browser, o.common.BrowserVersion = uaO.Browser()
		o.common.OS = uaO.OS()
		o.common.Platform = uaO.Platform()
		o.common.Mobile = uaO.Mobile()
		if o.common.Platform == "" && ua == "CLICKYAB" {
			o.common.Platform = "Android"
			o.common.OS = "Android"
			o.common.Mobile = true
			o.common.Browser = "AndroidSDK"
		}
		// TODO : lazy set?
		o.common.PlatformID = cyos.FindOsID(o.common.Platform)

		return o, nil
	}
}

// SetAlexa try to set Alexa flag if available
func SetAlexa(ua string, headers http.Header) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		// In go headers are not case sensitive and ok with _ and -
		if strings.Contains(ua, "Alexa") || headers.Get("ALEXATOOLBAR-ALX_NS_PH") != "" {
			o.common.Alexa = true
		}

		return o, nil
	}
}

// SetRequest try to set request in context, also all query params needed by the process
func SetRequest(host, method string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.common.Host = host
		o.common.Method = method
		o.common.MegaImp = <-random.ID
		return o, nil
	}
}

// SetSchema try to find schema of the request based on the request headers
func SetSchema(r *http.Request) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.common.Scheme = "http"
		if r.TLS != nil {
			o.common.Scheme = "https"
		}
		if xh := strings.ToLower(r.Header.Get("X-Forwarded-Proto")); xh == "https" {
			o.common.Scheme = "https"
		}

		return o, nil
	}
}

// SetQueryParameters try to get query parameters from the request and set the
// proper field
func SetQueryParameters(u *url.URL, ref string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.common.Parent = u.Query().Get("loc")
		o.common.Referrer = u.Query().Get("ref")

		if o.common.Referrer == "" {
			o.common.Referrer = ref
		}
		return o, nil
	}
}

// SetMinCPC try to set minimum cpc for this request
func SetMinCPC(i int64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.rtb.MinCPC = i
		return o, nil
	}
}

// SetMinBidPercentage try to set minimum bid for this request (normally from Website/App data)
func SetMinBidPercentage(i float64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.rtb.MinBidPercentage = i
		return o, nil
	}
}

// SetFloorDiv try to set floor div (the real floor is the floor/this value)
func SetFloorDiv(i int64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if i == 0 {
			i = 1
		}
		o.rtb.FloorDIV = i
		return o, nil
	}
}

// SetAllowUnderFloor try to set floor div (the real floor is the floor/this value)
func SetAllowUnderFloor(u bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.rtb.UnderFloor = u
		return o, nil
	}
}

// SetApp set application publisher
func SetApp(app *models.App) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.data.Website != nil || o.data.App != nil {
			return nil, fmt.Errorf("Website/App is already set")
		}
		o.data.App = app
		return o, nil
	}
}

// SetWebsite is the webiste publisher
func SetWebsite(web *models.Website) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.data.Website != nil || o.data.App != nil {
			return nil, fmt.Errorf("Website/App is already set")
		}
		o.data.Website = web
		return o, nil
	}
}

// SetEventPage in request that need it
func SetEventPage(ep string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.rtb.EventPage = ep
		return o, nil
	}
}

// SetAsync make this request an async request (vast mainly) default is sync
func SetAsync() ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.rtb.Async = true
		return o, nil
	}
}

// SetNoCap make this request to not use capping system
func SetNoCap() ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.rtb.NoCap = true
		return o, nil
	}
}

// SetNoTiny is the option to remove the tiny clickyab marker
func SetNoTiny() ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.common.NoTiny = true
		return o, nil
	}
}
