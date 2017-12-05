package builder

import (
	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models"
)

// Context is the app Context
type Context struct {
	typ string // app,native,vast,app

	ip             net.IP
	ua             string
	os             models.OS
	location       entity.Location
	browser        string
	browserVersion string

	tid       string
	publisher entity.Publisher
	seats     []entity.Seat
	cat       []entity.Category
	protocol  entity.Protocol
	user      entity.User

	NoTiny bool

	host     string
	method   string
	referrer string
	parent   string

	currencyRate float64
	noCap        bool
	noTiny       bool
	noShowT      bool

	floorDiv int64
}

func (c *Context) Protocol() entity.Protocol {
	return c.protocol
}

func (c *Context) User() entity.User {
	return c.user
}

func (c *Context) IP() net.IP {
	return c.ip
}

func (c *Context) OS() entity.OS {
	return c.os
}

func (c *Context) ClientID() string {
	return c.cid
}

func (c *Context) UserAgent() string {
	return c.ua
}

func (c *Context) Location() entity.Location {
	return c.location
}

func (c *Context) TrackID() string {
	return c.tid
}

func (c *Context) Publisher() entity.Publisher {
	return c.publisher
}

func (c *Context) Seats() []entity.Seat {
	return c.seats
}

func (c *Context) Category() []entity.Category {
	return c.cat
}
