package builder

import (
	"net"

	"time"

	"clickyab.com/crane/demand/entity"
)

// Context is the app Context
type Context struct {
	ts  time.Time
	typ entity.RequestType

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
	referrer string
	parent   string

	noCap            bool
	noTiny           bool
	multiVideo       bool
	floorPercentage  int64
	floorCPM         int64
	softFloorCPM     int64
	bidType          entity.BIDType
	minBidPercentage int64

	noShowT bool // Default is showT frame. but for exchange may be we don't want that

	suspicious int
	rate       float64

	networkID,
	brandID,
	carrierID int64

	networkName,
	brandName,
	carrierName string
}

func (c *Context) Network() (string, int64) {
	return c.networkName, c.networkID
}

func (c *Context) Carrier() (string, int64) {
	return c.carrierName, c.carrierID
}

func (c *Context) Brand() (string, int64) {
	return c.brandName, c.brandID
}

// Rate return the rate of this request
func (c *Context) Rate() float64 {
	return c.rate
}

// FloorCPM is the minimum CPM for this request, default is MinBid * DefaultCTR
// Currency is always Rial
func (c *Context) FloorCPM() int64 {
	return int64(float64(c.floorCPM)/100) * c.FloorPercentage()
}

// SoftFloorCPM is the minimum CPM for this request, default is FloorCPM
// Currency is always Rial
func (c *Context) SoftFloorCPM() int64 {
	return int64(float64(c.softFloorCPM)/100) * c.FloorPercentage()
}

// BIDType is the Bid Type for this request, no matter what, the system use
// cpc , but this is helpful for better data in logs/impression/click
func (c *Context) BIDType() entity.BIDType {
	return c.bidType
}

// MinBIDPercentage return the percentage for min bid
func (c *Context) MinBIDPercentage() int64 {
	if c.minBidPercentage <= 0 {
		c.minBidPercentage = 100
	}
	if c.minBidPercentage > 200 {
		c.minBidPercentage = 200
	}

	return c.minBidPercentage
}

// Suspicious return zero on ok status and a number on invalid value
func (c *Context) Suspicious() int {
	return c.suspicious
}

// Timestamp return the timestamp of the request
func (c *Context) Timestamp() time.Time {
	return c.ts
}

// Type return the request type
func (c *Context) Type() entity.RequestType {
	return c.typ
}

// Referrer is the request referrer
func (c *Context) Referrer() string {
	return c.referrer
}

// Parent is the request parent
func (c *Context) Parent() string {
	return c.parent
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

// FloorPercentage return the floor percentage for this request
// TODO : multiple floor in database
func (c *Context) FloorPercentage() int64 {
	if c.floorPercentage <= 0 {
		c.floorPercentage = 100
	}

	if c.floorPercentage > 200 {
		c.floorPercentage = 200
	}

	return c.floorPercentage
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
