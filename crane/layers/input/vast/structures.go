package vast

import (
	"net"

	"clickyab.com/crane/crane/entity"
)

type impression struct {
	trackID    string
	clientID   string
	ip         net.IP
	userAgent  string
	publisher  entity.Publisher
	location   entity.Location
	os         entity.OS
	slots      []entity.Seat
	categories []entity.Category
	attributes map[string]string
	protocol   string
}

func (m impression) IP() net.IP {
	return m.ip
}

func (m impression) OS() entity.OS {
	return m.os
}

func (m impression) ClientID() string {
	return m.clientID
}

func (m impression) Protocol() string {
	return m.protocol
}

func (m impression) UserAgent() string {
	return m.userAgent
}

func (m impression) Location() entity.Location {
	return m.location
}

func (m impression) Attributes() map[string]string {
	return m.attributes
}

func (m impression) TrackID() string {
	return m.trackID
}

func (m impression) Publisher() entity.Publisher {
	return m.publisher
}

func (m impression) Slots() []entity.Slot {
	return m.slots
}

func (m impression) Category() []entity.Category {
	return m.categories
}
