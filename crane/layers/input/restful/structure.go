package restful

import (
	"net/http"

	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
)

type restInput struct {
	payload
	request
}

type request struct {
	r *http.Request
}

type payload struct {
	FTrackID    string `json:"track_id"`
	FIP         net.IP `json:"ip"`
	FUserAgent  string `json:"user_agent"`
	Scheme      string `json:"scheme"`
	PageTrackID string `json:"page_track_id"`
	UserTrackID string `json:"user_track_id"`
	// Source return the publisher that this client is going into system from that
	FSource *local.Publisher `json:"source"`
	// Location of the request
	Locationn *local.Location `json:"location"`
	// Attributes is the generic attribute system
	FAttributes map[string]interface{} `json:"attributes"`
	// Slots is the slot for this request
	FSlots []*local.Slot `json:"slots"`
	// Category returns category obviously
	FCategory []entity.Category `json:"category"`
	// Platform return the publisher Platform
	Platform string `json:"platform"`
	// Is this publisher accept under floor ads or not ?
	UnderFloor *bool `json:"under_floor"`

	// generated field
	os entity.OS
}
