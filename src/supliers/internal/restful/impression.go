package restful

import (
	"entity"
	"net"
	"net/http"
	"services/random"
)

type impressionRest struct {
	r             *http.Request
	SIP           string                `json:"ip"`
	Mega          string                `json:"track_id"`
	UA            string                `json:"user_agent"`
	Pub           *restPublisher        `json:"source"`
	Loc           entity.Location       `json:"location"`
	ImpOS         entity.OS             `json:"os"`
	ImpSlots      []*slotRest           `json:"slots"`
	Categories    []entity.Category     `json:"categories"`
	ImpType       entity.ImpressionType `json:"type"`
	UnderFloorCPM bool                  `json:"under_floor"`

	dum []entity.Slot
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

func (ir impressionRest) OS() entity.OS {
	return ir.ImpOS
}

func (impressionRest) Attributes(entity.ImpressionAttributes) interface{} {
	return nil
}

func (ir impressionRest) Slots() []entity.Slot {
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

func newImpressionFromAppRequest(sup entity.Supplier, r *requestBody) (entity.Impression, error) {
	resp := impressionRest{}
	return &resp, nil
}

func newImpressionFromVastRequest(sup entity.Supplier, r *requestBody) (entity.Impression, error) {
	resp := impressionRest{}
	return &resp, nil
}

func newImpressionFromWebRequest(sup entity.Supplier, r *requestBody) (entity.Impression, error) {
	resp := impressionRest{}
	return &resp, nil
}
