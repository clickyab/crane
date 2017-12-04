package native

import (
	"net"

	"clickyab.com/crane/crane/entity"
)

type impression struct {
	attr       map[string]string
	trackID    string
	clientID   string
	ip         net.IP
	ua         string
	pub        entity.Publisher
	location   entity.Location
	os         entity.OS
	slots      []entity.Seat
	categories []entity.Category
	protocol   string
}

func (i *impression) TrackID() string {
	return i.trackID
}

func (i *impression) ClientID() string {
	return i.clientID
}

func (i *impression) IP() net.IP {
	return i.ip
}

func (i *impression) UserAgent() string {
	return i.ua
}
func (i *impression) Location() entity.Location {
	return i.location
}

func (i *impression) OS() entity.OS {
	return i.os
}

func (i *impression) Slots() []entity.Seat {
	return i.slots
}

func (i *impression) Category() []entity.Category {
	return i.categories
}

func (i *impression) Publisher() entity.Publisher {
	return i.pub
}

func (i *impression) Protocol() string {
	return i.protocol
}
