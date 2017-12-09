package builder

import (
	"fmt"
	"net"
	"net/url"

	"net/http"

	"strings"

	"clickyab.com/crane/crane/builder/internal/cyos"
	"clickyab.com/crane/crane/models"

	"crypto/sha1"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/mssola/user_agent"
)

// SetType is the type setter for context
func SetType(typ string) ShowOptionSetter {
	return func(options *Context) (*Context, error) {
		if typ != "vast" && typ != "App" && typ != "web" && typ != "native" {
			return nil, fmt.Errorf("type is not supported %s", typ)
		}
		options.typ = typ
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

// SetIPLocation is the IP and location setter for context, also it extract the IP information
func SetIPLocation(ip string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		ipv4 := net.ParseIP(ip)
		if ipv4 == nil {
			return nil, fmt.Errorf("invalid IP %s", ip)
		}
		o.ip = ipv4
		l := models.GetProvinceISPByIP(ipv4)
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
		}
		return o, nil
	}
}

// SetRequest try to set request in context, also all query params needed by the process
func SetRequest(host, method string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.host = host
		o.method = method
		return o, nil
	}
}

// SetProtocol try to find protocol of the request based on the request headers
func SetProtocol(r *http.Request) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
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

// SetEventPage is the event page setter for multi request
func SetEventPage(ep string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.eventPage = ep
		return o, nil
	}
}

// SetDisableCapping disable capping?
func SetDisableCapping(disable bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.noCap = disable
		return o, nil
	}
}

var copLen = config.RegisterInt("clickyab.cop_len", 10, "cop key len")

// SetTID try to set tid
func SetTID(id string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.ua == "" || o.ip == nil {
			return nil, fmt.Errorf("use this after setting ip and ua")
		}
		o.tid = id
		if o.tid == "" {
			o.tid = createHash(copLen.Int(), []byte(o.ip), []byte(o.ua))
		}

		o.user = user(o.tid)
		return o, nil
	}
}

// SetAlexa try to set Alexa flag if available
func SetAlexa(ua string, headers http.Header) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		// In go headers are not case sensitive and ok with _ and -
		if strings.Contains(ua, "Alexa") || headers.Get("ALEXATOOLBAR-ALX_NS_PH") != "" {
			o.alexa = true
		}

		return o, nil
	}
}

// SetQueryParameters try to get query parameters from the request and
// set the proper field
func SetQueryParameters(u *url.URL, ref string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.parent = u.Query().Get("parent")
		o.referrer = u.Query().Get("ref")
		if o.referrer == "" {
			o.referrer = ref
		}
		return o, nil
	}
}

// SetCurrencyRate set currency convert rate to rial
func SetCurrencyRate(a float64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.currencyRate = a
		return o, nil
	}
}

// SetFloorDiv try to set floor div (the real floor is the floor/this value)
func SetFloorDiv(i int64) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if i == 0 {
			i = 1
		}
		o.floorDiv = i
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

// SetNoTiny is the option to remove the tiny clickyab marker
func SetNoTiny(noTiny bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.noTiny = noTiny
		return o, nil
	}
}

// SetMultiVideo is the option to remove the tiny clickyab marker
func SetMultiVideo(v bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		o.multiVideo = v
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
