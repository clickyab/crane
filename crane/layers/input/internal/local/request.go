package local

import (
	"net"

	"clickyab.com/crane/crane/entity"
)

type request struct {
	attr      map[string]string
	ip        net.IP
	os        entity.OS
	client    string
	protocol  string
	userAgent string
	location  entity.Location
}

func (r *request) IP() net.IP {
	return r.ip
}

func (r *request) OS() entity.OS {
	return r.os
}

func (r *request) ClientID() string {
	return r.client
}

func (r *request) Protocol() string {
	return r.protocol
}

func (r *request) UserAgent() string {
	return r.userAgent
}

func (r *request) Location() entity.Location {
	return r.location
}

func (r *request) Attributes() map[string]string {
	return r.attr
}
