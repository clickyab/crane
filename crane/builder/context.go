package builder

import (
	"net"

	"clickyab.com/crane/crane/entity"
)

// Context is the app Context
type Context struct {
	ip        net.IP
	os        entity.OS
	cid       string
	ua        string
	location  entity.Location
	tid       string
	publisher entity.Publisher
	seats     []entity.Seat
	cat       []entity.Category
	protocol  entity.Protocol
	user      entity.User
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
