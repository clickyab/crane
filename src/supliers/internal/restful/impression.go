package restful

import (
	"entity"
	"net"
	"services/gmaps"
	"services/ip2location/client"
	"services/random"
)

type impressionRest struct {
	SIP           string                `json:"ip"`
	Mega          string                `json:"track_id"`
	UA            string                `json:"user_agent"`
	Pub           *restPublisher        `json:"source"`
	Loc           entity.Location       `json:"location"`
	ImpSlots      []*slotRest           `json:"slots"`
	Categories    []entity.Category     `json:"categories"`
	ImpType       entity.ImpressionType `json:"type"`
	UnderFloorCPM bool                  `json:"under_floor"`

	Attr map[string]interface{} `json:"attributes"`

	dum    []entity.Slot
	latlon entity.LatLon
}

type location struct {
	TheCountry  entity.Country  `json:"country"`
	TheProvince entity.Province `json:"province"`
	TheLatLon   entity.LatLon   `json:"latlon"`
}

func (l location) Country() entity.Country {
	return l.TheCountry
}

func (l location) Province() entity.Province {
	return l.TheProvince
}

func (l location) LatLon() entity.LatLon {
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

func (ir impressionRest) Source() entity.Publisher {
	return ir.Pub
}

func (ir impressionRest) Location() entity.Location {
	return ir.Loc
}

func (impressionRest) Attributes() map[string]interface{} {
	return nil
}

func (ir *impressionRest) Slots() []entity.Slot {
	if ir.dum == nil {
		ir.dum = make([]entity.Slot, len(ir.ImpSlots))
		for i := range ir.ImpSlots {
			ir.dum[i] = ir.ImpSlots[i]
		}
	}
	return ir.dum
}

func (ir impressionRest) Category() []entity.Category {
	return ir.Categories
}

func (ir impressionRest) Type() entity.ImpressionType {
	return ir.ImpType
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
		TheCountry: entity.Country{
			Name:  d.CountryLong,
			ISO:   d.CountryShort,
			Valid: d.CountryLong != "-",
		},

		TheProvince: entity.Province{
			Valid: d.Region != "-",
			Name:  d.Region,
		},

		TheLatLon: ir.latlon,
	}

}

func newImpressionFromAppRequest(sup entity.Supplier, r *requestBody) (*impressionRest, error) {
	resp := impressionRest{
		SIP:           r.IP,
		UA:            r.App.UserAgent,
		ImpType:       entity.ImpressionTypeApp,
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
	resp.latlon = entity.LatLon{
		Valid: err == nil,
		Lat:   lat,
		Lon:   lon,
	}
	resp.Pub.sup = sup
	resp.extractData()
	return &resp, nil
}

func newImpressionFromVastRequest(sup entity.Supplier, r *requestBody) (*impressionRest, error) {
	resp := impressionRest{
		SIP:           r.IP,
		UA:            r.Vast.UserAgent,
		ImpType:       entity.ImpressionTypeVast,
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

func newImpressionFromWebRequest(sup entity.Supplier, r *requestBody) (*impressionRest, error) {
	resp := impressionRest{
		SIP:           r.IP,
		UA:            r.Web.UserAgent,
		ImpType:       entity.ImpressionTypeWeb,
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
