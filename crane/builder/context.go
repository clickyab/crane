package builder

import (
	"net"

	"clickyab.com/crane/crane/entity"
)

// Context is the app Context
type Context struct {
	typ string // app,native,vast,app

	ip             net.IP
	ua             string
	os             entity.OS
	location       entity.Location
	browser        string
	browserVersion string

	tid       string
	publisher entity.Publisher
	seats     []entity.Seat
	cat       []entity.Category
	protocol  entity.Protocol
	user      entity.User
	currency  string

	eventPage string
	alexa     bool
	NoTiny    bool

	host     string
	method   string
	referrer string
	parent   string

	currencyRate float64
	noCap        bool
	noTiny       bool
	multiVideo   bool

	floorDiv int64
}

func (c *Context) Alexa() bool {
	return c.alexa
}

func (c *Context) EventPage() string {
	return c.eventPage
}

func (c *Context) Capping() bool {
	return !c.noCap
}

func (c *Context) IsMobile() bool {
	return c.os.Mobile
}

func (c *Context) Isp() string {
	return c.location.ISP().Name
}

func (c *Context) FloorDiv() int64 {
	return c.floorDiv
}

func (c *Context) Tiny() bool {
	return !c.noTiny
}

func (c *Context) Currency() string {
	return c.currency
}

func (c *Context) MultiVideo() bool {
	return c.multiVideo
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
	return c.tid
}

func (c *Context) UserAgent() string {
	return c.ua
}

func (c *Context) Location() entity.Location {
	return c.location
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
