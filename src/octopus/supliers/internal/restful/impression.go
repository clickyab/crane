package restful

import (
	"net"
	"octopus/exchange"
	"services/gmaps"
	"services/ip2location/client"
	"services/random"
)

type impressionRest struct {
	SIP           string                      `json:"ip"`
	Mega          string                      `json:"track_id"`
	UA            string                      `json:"user_agent"`
	Pub           *restPublisher              `json:"source"`
	Loc           exchange.Location           `json:"location"`
	ImpSlots      []*slotRest                 `json:"slots"`
	Categories    []exchange.Category         `json:"categories"`
	ImpPlatform   exchange.ImpressionPlatform `json:"platform"`
	UnderFloorCPM bool                        `json:"under_floor"`

	Attr map[string]interface{} `json:"attributes"`

	dum    []exchange.Slot
	latlon exchange.LatLon
}

type location struct {
	TheCountry  exchange.Country  `json:"country"`
	TheProvince exchange.Province `json:"province"`
	TheLatLon   exchange.LatLon   `json:"latlon"`
}

func (l location) Country() exchange.Country {
	return l.TheCountry
}

func (l location) Province() exchange.Province {
	return l.TheProvince
}

func (l location) LatLon() exchange.LatLon {
	return l.TheLatLon
}

func (ir *impressionRest) TrackID() string {
	if ir.Mega == "" {
		ir.Mega = <-random.ID
	}

	return ir.Mega
}

func (ir impressionRest) IP() net.IP {
	return net.ParseIP(ir.SIP)
}

func (ir impressionRest) UserAgent() string {
	return ir.UA
}

func (ir impressionRest) Source() exchange.Publisher {
	return ir.Pub
}

func (ir impressionRest) Location() exchange.Location {
	return ir.Loc
}

func (impressionRest) Attributes() map[string]interface{} {
	return nil
}

func (ir *impressionRest) Slots() []exchange.Slot {
	if ir.dum == nil {
		ir.dum = make([]exchange.Slot, len(ir.ImpSlots))
		for i := range ir.ImpSlots {
			ir.dum[i] = ir.ImpSlots[i]
		}
	}
	return ir.dum
}

func (ir impressionRest) Category() []exchange.Category {
	return ir.Categories
}

func (ir impressionRest) Platform() exchange.ImpressionPlatform {
	return ir.ImpPlatform
}

func (ir impressionRest) UnderFloor() bool {
	return ir.UnderFloorCPM
}

func (ir impressionRest) Raw() interface{} {
	return ir
}

func (ir *impressionRest) extractData() {
	d := client.IP2Location(ir.SIP)
	//logrus.Debug(d)
	ir.Loc = location{
		TheCountry: exchange.Country{
			Name:  d.CountryLong,
			ISO:   d.CountryShort,
			Valid: d.CountryLong != "-",
		},

		TheProvince: exchange.Province{
			Valid: d.Region != "-",
			Name:  d.Region,
		},

		TheLatLon: ir.latlon,
	}

}

func newImpressionFromAppRequest(sup exchange.Supplier, r *requestBody) (*impressionRest, error) {
	resp := impressionRest{
		SIP:           r.IP,
		UA:            r.App.UserAgent,
		ImpPlatform:   exchange.ImpressionPlatformApp,
		Categories:    r.Categories,
		ImpSlots:      r.Slots,
		Mega:          <-random.ID,
		UnderFloorCPM: r.UnderFloor,
		Pub:           r.Publisher,
		Attr: map[string]interface{}{
			"network":     r.App.Network,
			"brand":       r.App.Brand,
			"cid":         r.App.CID,
			"lac":         r.App.LAC,
			"mcc":         r.App.MCC,
			"mnc":         r.App.MNC,
			"language":    r.App.Language,
			"model":       r.App.Model,
			"operator":    r.App.Operator,
			"os_identity": r.App.OSIdentity,
		},
	}
	lat, lon, err := gmaps.LockUp(r.App.MCC, r.App.MNC, r.App.LAC, r.App.CID)
	resp.latlon = exchange.LatLon{
		Valid: err == nil,
		Lat:   lat,
		Lon:   lon,
	}
	resp.Pub.sup = sup
	resp.extractData()
	return &resp, nil
}

func newImpressionFromVastRequest(sup exchange.Supplier, r *requestBody) (*impressionRest, error) {
	resp := impressionRest{
		SIP:           r.IP,
		UA:            r.Vast.UserAgent,
		ImpPlatform:   exchange.ImpressionPlatformVast,
		Categories:    r.Categories,
		ImpSlots:      r.Slots,
		Mega:          <-random.ID,
		UnderFloorCPM: r.UnderFloor,

		Attr: map[string]interface{}{
			"referrer": r.Vast.Referrer,
			"parent":   r.Vast.Parent,
		},
		Pub: r.Publisher,
	}
	resp.Pub.sup = sup
	resp.extractData()
	return &resp, nil
}

func newImpressionFromWebRequest(sup exchange.Supplier, r *requestBody) (*impressionRest, error) {
	resp := impressionRest{
		SIP:           r.IP,
		UA:            r.Web.UserAgent,
		ImpPlatform:   exchange.ImpressionPlatformWeb,
		Categories:    r.Categories,
		ImpSlots:      r.Slots,
		Mega:          <-random.ID,
		UnderFloorCPM: r.UnderFloor,

		Attr: map[string]interface{}{
			"referrer": r.Web.Referrer,
			"parent":   r.Web.Parent,
		},
		Pub: r.Publisher,
	}
	resp.Pub.sup = sup
	resp.extractData()
	return &resp, nil
}
