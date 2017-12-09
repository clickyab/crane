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

// Alexa return if user has alexa installed
func (c *Context) Alexa() bool {
	return c.alexa
}

// EventPage is the event page for the single ad requests
func (c *Context) EventPage() string {
	return c.eventPage
}

// Capping is enabled or not?
func (c *Context) Capping() bool {
	return !c.noCap
}

// IsMobile return if the request is from mobile
func (c *Context) IsMobile() bool {
	return c.os.Mobile
}

// ISP return the isp of this request
func (c *Context) ISP() string {
	return c.location.ISP().Name
}

// FloorDiv return the floor div for this request
// TODO : multiple floor in database
func (c *Context) FloorDiv() int64 {
	return c.floorDiv
}

// Tiny means we need to show the tiny mark in ad
func (c *Context) Tiny() bool {
	return !c.noTiny
}

// Currency is for this request currency
func (c *Context) Currency() string {
	return c.currency
}

// MultiVideo means this request can be multi video
func (c *Context) MultiVideo() bool {
	return c.multiVideo
}

// Protocol return the protocol of this request (http/https)
func (c *Context) Protocol() entity.Protocol {
	return c.protocol
}

// User return the current user of this request based on finger print
func (c *Context) User() entity.User {
	return c.user
}

// IP return the ip of this request
func (c *Context) IP() net.IP {
	return c.ip
}

// OS of this request (UA based)
func (c *Context) OS() entity.OS {
	return c.os
}

// UserAgent return the user agent
func (c *Context) UserAgent() string {
	return c.ua
}

// Location return the location based on ip 2 location database
func (c *Context) Location() entity.Location {
	return c.location
}

// Publisher return the current publisher
func (c *Context) Publisher() entity.Publisher {
	return c.publisher
}

// Seats return seats (slots)
func (c *Context) Seats() []entity.Seat {
	return c.seats
}

// Category return the request categories
func (c *Context) Category() []entity.Category {
	return c.cat
}
