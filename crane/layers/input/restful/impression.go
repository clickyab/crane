package restful

import (
	"net"

	"clickyab.com/crane/crane/entity"
	"github.com/mssola/user_agent"
)

type impression struct {
	FTrackID    string `json:"track_id"`
	FIP         net.IP `json:"ip"`
	FUserAgent  string `json:"user_agent"`
	FProtocol   string `json:"scheme"`
	UserTrackID string `json:"user_track_id"`
	// FPublisher return the publisher that this client is going into system from that
	// IF YOU GET ERROR IT'S POSSIBLE THAT JSON TAG SHOULD BE RENAME TO source
	FPublisher entity.Publisher `json:"publisher"`
	// Location of the request
	FLocation entity.Location `json:"location"`
	// Attributes is the generic attribute system
	FAttributes map[string]string `json:"attributes"`
	// Slots is the slot for this request
	FSlots []entity.Seat `json:"slots"`
	// Category returns category obviously
	FCategory []entity.Category `json:"category"`
	// Platform return the publisher Platform
	Platform string `json:"platform"`
	// Is this publisher accept under floor ads or not ?
	UnderFloor *bool `json:"under_floor"`
}

func (r impression) TrackID() string {
	return r.FTrackID
}

func (r impression) Publisher() entity.Publisher {
	return r.FPublisher
}

func (r impression) Slots() []entity.Seat {
	return r.FSlots
}

func (r impression) Category() []entity.Category {
	return r.FCategory
}

func (r impression) IP() net.IP {
	return r.FIP
}

func (r impression) OS() entity.OS {
	u := user_agent.New(r.FUserAgent)
	return entity.OS{
		Valid:  u.OS() != "",
		Mobile: u.Mobile(),
		Name:   u.OS(),
	}
}

func (r impression) ClientID() string {
	return r.UserTrackID
}

func (r impression) Protocol() string {
	return r.FProtocol
}

func (r impression) UserAgent() string {
	return r.FUserAgent
}

func (r impression) Location() entity.Location {
	return r.FLocation
}

func (r impression) Attributes() map[string]string {
	return r.FAttributes
}
