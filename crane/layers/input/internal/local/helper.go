package local

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net"
	"net/http"
	"strings"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/ip2location"
	"github.com/clickyab/services/random"
	"github.com/mssola/user_agent"
	"github.com/rs/xmux"
)

const (
	httpsScheme = "https"
	httpScheme  = "http"
)

// Request fill the request interface for every get requests
// extra can be request body (json for example or anything that may needed later
// the flow is something like: route -> Request -> Target (Vast, Native, etc) -> ...
func Request(ctx context.Context, r *http.Request, extra map[string]string) entity.Request {
	ip := net.ParseIP(framework.RealIP(r))
	clientUserLocation := location(ip)

	os := os(r.UserAgent())
	protocol := httpScheme
	if r.TLS != nil {
		protocol = httpsScheme
	}
	if xh := strings.ToLower(r.Header.Get("X-Forwarded-Proto")); xh == httpsScheme {
		protocol = httpsScheme
	}

	var attr = make(map[string]string)
	for k, v := range r.URL.Query() {
		attr[k] = v[0]
	}

	for _, p := range xmux.Params(ctx) {
		attr[p.Name] = p.Value
	}
	for k, v := range extra {
		attr[k] = v
	}

	re := &request{
		attr:      attr,
		ip:        ip,
		os:        os,
		userAgent: r.UserAgent(),
		location:  clientUserLocation,
		protocol:  protocol,
	}
	re.client = cidGenerator(re)
	return re
}

// ExtractSlot ExtractSlot
func ExtractSlot(supplier, publisher string, typ entity.Platforms, width, height int, uniqueID string, attributes map[string]interface{}) entity.Slot {
	return &Slot{
		FID:       fmt.Sprintf("%s/%s/%s/%dx%d/%s", supplier, publisher, typ, width, height, uniqueID),
		attribute: attributes,
		FHeight:   height,
		FWidth:    width,
		FTrackID:  <-random.ID,
		// TODO read from config
		slotCTR: 0.1,
	}
}

func os(ua string) entity.OS {
	userAgent := user_agent.New(ua)
	return entity.OS{
		Valid:  userAgent.OS() != "",
		Name:   userAgent.OS(),
		Mobile: userAgent.Mobile(),
	}

}

func location(ip net.IP) Location {
	l := Location{}
	if ip.String() == "" {
		return l
	}
	ip2loc := ip2location.IP2Location(ip.String())

	return Location{
		FCountry: entity.Country{
			Name:  ip2loc.CountryLong,
			ISO:   ip2loc.CountryShort,
			Valid: ip2loc.CountryShort != "",
		},
		FProvince: entity.Province{
			Name:  ip2loc.Region,
			Valid: ip2loc.Region != "",
		},
		FLatLon: entity.LatLon{
			Valid: false,
		},
	}
}

var (
	copLen = config.RegisterInt("crane.cop.len", 10, "client key len")
)

func cidGenerator(r entity.Request) string {
	tID := r.Attributes()["tid"]
	if len(tID) < 10 {
		return createHash(copLen.Int(), []byte(r.UserAgent()), []byte(r.IP()))
	}
	return tID
}

func createHash(l int, items ...[]byte) string {
	h := sha1.New()
	for i := range items {
		_, _ = h.Write(items[i])
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
