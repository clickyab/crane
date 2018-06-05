package builder

import (
	"errors"
	"fmt"
	"net"

	"net/http"

	"strings"

	"clickyab.com/crane/internal/cyos"

	"crypto/sha1"

	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/cell"
	"clickyab.com/crane/models/ip2l"
	"github.com/bsm/openrtb"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/mssola/user_agent"
)

// SetRate set the timestamp. must be first!
func SetRate(r float64) ShowOptionSetter {
	return func(options *Context) (*Context, error) {
		options.rate = r
		return options, nil
	}
}

// SetTimestamp set the timestamp. must be first!
func SetTimestamp() ShowOptionSetter {
	return func(options *Context) (*Context, error) {
		options.ts = time.Now()
		return options, nil
	}
}

// SetCurrency is the type setter for context
func SetCurrency(c string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if c != "IRR" && c != "USD" {
			return o, fmt.Errorf("%s is not valid currency", c)
		}
		o.currency = c
		return o, nil
	}
}

// SetUnderfloor is the type setter for context
func SetUnderfloor(c bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.underfloor = c
		return o, nil
	}
}

// SetIPLocation is the IP and location setter for context, also it extract the IP information
func SetIPLocation(ip string, user *openrtb.User, device *openrtb.Device) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ipv4 := net.ParseIP(ip)
		if ipv4 == nil {
			return nil, fmt.Errorf("invalid IP %s", ip)
		}
		o.ip = ipv4
		var lat, lon float64
		if device != nil && device.Geo != nil {
			lat, lon = device.Geo.Lat, device.Geo.Lon

		}
		if lat == 0 && lon == 0 {
			if user != nil && user.Geo != nil {
				lat, lon = user.Geo.Lat, user.Geo.Lon
			}
		}

		l := ip2l.GetProvinceISPByIP(ipv4, lat, lon)
		o.location = l
		return o, nil
	}
}

// SetOSUserAgent try to set user agent and os and all related things
func SetOSUserAgent(ua string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		uaO := user_agent.New(ua)
		o.ua = ua
		os := uaO.OS()
		platForm := uaO.Platform()
		mobile := uaO.Mobile()
		o.browser, o.browserVersion = uaO.Browser()
		osID := cyos.FindOsID(platForm)
		o.os = entity.OS{
			Name:   os,
			Valid:  osID != 0,
			Mobile: mobile,
			ID:     osID,
		}
		return o, nil
	}
}

// SetTargetHost try to set request in context, also all query params needed by the process
func SetTargetHost(host string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.host = host
		return o, nil
	}
}

// SetTrueView set true if true view
func SetTrueView(tv bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.tv = tv
		return o, nil
	}
}

// SetProtocolByRequest try to find protocol of the request based on the request headers
func SetProtocolByRequest(r *http.Request) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		// TODO : Create framework.Schema() function
		o.protocol = entity.HTTP
		if r.TLS != nil {
			o.protocol = entity.HTTPS
		}
		if xh := strings.ToLower(r.Header.Get("X-Forwarded-Proto")); xh == "https" {
			o.protocol = entity.HTTPS
		}

		return o, nil
	}
}

// SetProtocol try to find protocol of the request based on function parameters
func SetProtocol(scheme entity.Protocol) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.protocol = scheme
		return o, nil
	}
}

// SetEventPage is the event page setter for multi request
func SetEventPage(ep string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.eventPage = ep
		return o, nil
	}
}

// SetCategory set the capping mode
func SetCategory(b *openrtb.BidRequest) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		var category []entity.Category
		if b.Site != nil {
			for _, v := range b.Site.Cat {
				if len(v) > 3 {
					category = append(category, entity.Category(v[3:]))
				}
			}
		} else if b.App != nil {
			for _, v := range b.App.Cat {
				if len(v) > 3 {
					category = append(category, entity.Category(v[3:]))
				}
			}
		}
		o.cat = category
		return o, nil
	}
}

// SetCappingMode set the capping mode
func SetCappingMode(mode entity.CappingMode) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.cappingMode = mode
		return o, nil
	}
}

// SetFatFinger enable/disable fat finger? default is disable
func SetFatFinger(ff bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.fatFinger = ff
		return o, nil
	}
}

var copLen = config.RegisterInt("crane.context.cop_len", 10, "cop key len")

// SetTID try to set tid
func SetTID(id string, ip, ua string, extra ...string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.ua == "" || o.ip == nil {
			return nil, fmt.Errorf("use this after setting ip and ua")
		}
		ee := make([][]byte, len(extra)+3)
		ee[0] = []byte(id)
		ee[1] = []byte(ip)
		ee[2] = []byte(ua)
		for i := range extra {
			ee[i+3] = []byte(extra[i])
		}
		o.tid = createHash(copLen.Int(), ee...)

		o.user = user(o.tid)
		return o, nil
	}
}

// SetParent is same as SetQueryParameters just for setting parent for demands
func SetParent(page, ref string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.parent = page
		o.referrer = ref
		return o, nil
	}
}

// SetFloorPercentage try to set floor div (the real floor is the floor/this value)
func SetFloorPercentage(i int64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if i <= 0 {
			return nil, fmt.Errorf("invalid floor value")
		}
		o.floorPercentage = i
		return o, nil
	}
}

// SetMinBidPercentage set the minimum bid percentage
func SetMinBidPercentage(i int64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if i <= 0 {
			return nil, fmt.Errorf("invalid floor value")
		}
		o.minBidPercentage = i
		return o, nil
	}
}

// SetPublisher set publisher in context
func SetPublisher(pub entity.Publisher) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.publisher != nil {
			return nil, fmt.Errorf("publisher is already set")
		}
		o.publisher = pub
		return o, nil
	}
}

//SetSuspicious is the function to set suspicious code, default is zero
func SetSuspicious(suspCode int) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.suspicious = suspCode
		return o, nil
	}
}

// SetNoTiny is the option to remove the tiny clickyab marker
func SetNoTiny(noTiny bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.noTiny = noTiny
		return o, nil
	}
}

// SetStrategy is the option to set strategy for request
func SetStrategy(s []string, sup entity.Supplier) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if len(s) == 0 {
			o.strategy = sup.Strategy()
			return o, nil
		}
		st := entity.GetStrategy(s)
		if st.Valid() && st.IsSubsetOf(sup.Strategy()) {
			o.strategy = st
			return o, nil
		}
		return o, errors.New("not valid strategy")
	}
}

// SetMultiVideo is the option to remove the tiny clickyab marker
func SetMultiVideo(v bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.multiVideo = v
		return o, nil
	}
}

// SetBrand is set brand id from name
func SetBrand(v string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		n, err := cell.GetBrandByName(strings.ToLower(v))
		if err != nil {
			return o, err
		}
		o.brandName = n
		return o, nil
	}
}

// SetCarrier is set carrier id from name
func SetCarrier(v string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		n, err := cell.GetCarrierByName(v)
		if err != nil {
			return o, err
		}
		o.carrierName = n
		return o, nil
	}
}

// SetConnType is set network  2g,3g,...
func SetConnType(v int) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.connectionType = v
		return o, nil
	}
}

// SetPreventDefault is a simple workaround for our old sdk. the old sdk do not pass click to higher level
// and we use post message to handle that, pre 4 sdk are that way
func SetPreventDefault(prevent bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.preventDefault = prevent
		return o, nil
	}
}

// SetCreativesStatistics is set network creatives statistics
func SetCreativesStatistics(data []entity.CreativeStatistics) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.creativesStat = data
		return o, nil
	}
}

func createHash(l int, items ...[]byte) string {
	h := sha1.New()
	for i := range items {
		_, err := h.Write(items[i])
		assert.Nil(err)
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	if l >= len(sum) {
		l = len(sum)
	}
	if l < 1 {
		l = 1
	}
	return sum[:l]
}
