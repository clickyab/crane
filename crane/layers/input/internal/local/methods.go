package local

import (
	"fmt"
	"net"

	"context"
	"crypto/sha1"
	"net/http"

	"strings"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models/publisher"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/ip2location"
	"github.com/clickyab/services/random"
	"github.com/mssola/user_agent"
	"github.com/rs/xmux"
)

const (
	httpsScheme string = "https"
	httpScheme  string = "http"
)

// ExtractSlot ExtractSlot
func ExtractSlot(supplier, publisher string, typ publisher.Platforms, width, height int, uniqueID string, attributes map[string]interface{}) entity.Slot {
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

// OS OS
func OS(ua string) entity.OS {
	userAgent := user_agent.New(ua)
	return entity.OS{
		Valid:  userAgent.OS() != "",
		Name:   userAgent.OS(),
		Mobile: userAgent.Mobile(),
	}

}

// FLocation FLocation
func FLocation(ip net.IP) Location {
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

// CopCreator create client id
func CopCreator(ctx context.Context, r *http.Request) string {
	tID := r.URL.Query().Get("tid")
	if len(tID) < 10 {
		return createHash(copLen.Int(), []byte(r.UserAgent()), []byte(net.ParseIP(framework.RealIP(r))))
	}
	return tID
}

// createHash is used to handle the cop key
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

// GetRequestData fill the request interface for every get requests
func GetRequestData(ctx context.Context, r *http.Request) entity.Request {
	ip := net.ParseIP(framework.RealIP(r))
	clientUserLocation := FLocation(ip)
	clientID := CopCreator(ctx, r)
	os := OS(r.UserAgent())
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
	return &request{
		attr:     attr,
		ip:       ip,
		os:       os,
		location: clientUserLocation,
		client:   clientID,
		protocol: protocol,
	}
}
